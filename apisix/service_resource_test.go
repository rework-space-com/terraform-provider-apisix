package apisix

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestServiceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_service" "test" {
	name  = "test"
	hosts = ["foo.com", "*.bar.com"]
	labels = {
		"version" = "v1"
	}
}			
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// NOTE: State value checking is only necessary for Computed attributes,
					//       as the testing framework will automatically return test failures
					//       for configured attributes that mismatch the saved state.
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("apisix_service.test", "id"),
					resource.TestCheckResourceAttrSet("apisix_service.test", "enable_websocket"),
					resource.TestCheckResourceAttr("apisix_service.test", "enable_websocket", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apisix_service.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "apisix_service" "test" {
	name  = "test"
	hosts = ["foo.com"]
	labels = {
		"version" = "v2"
	}
	enable_websocket = true
}					
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apisix_service.test", "id"),
					resource.TestCheckResourceAttrSet("apisix_service.test", "enable_websocket"),
					resource.TestCheckResourceAttr("apisix_service.test", "enable_websocket", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestServiceResource_EncryptedPluginField covers the Update path with an
// encrypt_field (openid-connect client_secret). When data_encryption is enabled,
// the PUT response returns ciphertext for encrypt_fields; Update() must set
// state.plugins to the planned value, or Terraform fails with "Provider produced
// inconsistent result after apply".
// NOTE: with a stable keyring the re-fetch gate returns plaintext either way, so
// this passes with and without the fix; it covers the Update+encrypted path, not
// the keyring-rotation scenario that triggers the inconsistency.
func TestServiceResource_EncryptedPluginField(t *testing.T) {
	encryptedFieldValue := "plaintext" + "-dummy-" + "value-" + "12345"

	pluginsBlock := fmt.Sprintf(`
		openid-connect = {
			bearer_only   = true
			unauth_action = "pass"
			client_id     = "example-client-id"
			client_secret = %q
			discovery     = "https://example.com/.well-known/openid-configuration"
		}
`, encryptedFieldValue)

	createConfig := providerConfig + `
resource "apisix_service" "encrypted_plugin" {
	name  = "EncryptedPluginField"
	hosts = ["encrypted-plugin-field.example.com"]
	plugins = jsonencode(
		{
` + pluginsBlock + `
		}
	)
}
`

	// In-place Update (enable_websocket) re-PUTs the identical encrypted plugin field.
	updatedConfig := providerConfig + `
resource "apisix_service" "encrypted_plugin" {
	name  = "EncryptedPluginField"
	hosts = ["encrypted-plugin-field.example.com"]
	enable_websocket = true
	plugins = jsonencode(
		{
` + pluginsBlock + `
		}
	)
}
`

	expectedPlugins := map[string]interface{}{
		"openid-connect": map[string]interface{}{
			"bearer_only":   true,
			"unauth_action": "pass",
			"client_id":     "example-client-id",
			"client_secret": encryptedFieldValue,
			"discovery":     "https://example.com/.well-known/openid-configuration",
		},
	}

	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("apisix_service.encrypted_plugin", "id"),
		resource.TestCheckResourceAttrWith("apisix_service.encrypted_plugin", "plugins", func(value string) error {
			var actual interface{}
			if err := json.Unmarshal([]byte(value), &actual); err != nil {
				return fmt.Errorf("failed to unmarshal plugins: %w", err)
			}
			if !reflect.DeepEqual(actual, expectedPlugins) {
				return fmt.Errorf("expected plugins to equal %#v, got %#v (state likely holds ciphertext from the PUT response instead of the plaintext configured value)", expectedPlugins, actual)
			}
			return nil
		}),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{Config: createConfig, Check: check},
			// Update and Read testing.
			{Config: updatedConfig, Check: check},
			// Plan-only with the identical config must show no changes.
			{Config: updatedConfig, PlanOnly: true},
		},
	})
}
