package apisix

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPluginConfigResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_plugin_config" "test" {
	id   = "007"
	desc = "Example of the plugin config resource usage"
	plugins = jsonencode(
		{
			prometheus = {
				prefer_name = true
			}
		}
	)
	labels = {
		version = "v1"
	}
}	
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// NOTE: State value checking is only necessary for Computed attributes,
					//       as the testing framework will automatically return test failures
					//       for configured attributes that mismatch the saved state.
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("apisix_plugin_config.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apisix_plugin_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "apisix_plugin_config" "test" {
	id   = "007"
	desc = "Example of the plugin config resource usage"
	plugins = jsonencode(
		{
			prometheus = {
				prefer_name = false
			}
		}
	)
	labels = {
		version = "v2"
		env			= "stage"
	}
}				
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apisix_plugin_config.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestPluginConfigResource_EncryptedPluginField covers the Create + Update path with
// an encrypt_field (openid-connect client_secret) on plugin_config. Same shape and
// stable-keyring caveat as TestServiceResource_EncryptedPluginField: passes with and
// without the fix on a healthy APISIX; covers the Update+encrypted path, not the
// keyring-rotation scenario.
func TestPluginConfigResource_EncryptedPluginField(t *testing.T) {
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
resource "apisix_plugin_config" "encrypted_plugin" {
	id   = "enc-1"
	desc = "EncryptedPluginField v1"
	plugins = jsonencode(
		{
` + pluginsBlock + `
		}
	)
}
`

	// In-place Update (desc change) re-PUTs the identical encrypted plugin field.
	updatedConfig := providerConfig + `
resource "apisix_plugin_config" "encrypted_plugin" {
	id   = "enc-1"
	desc = "EncryptedPluginField v2"
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
		resource.TestCheckResourceAttrSet("apisix_plugin_config.encrypted_plugin", "id"),
		resource.TestCheckResourceAttrWith("apisix_plugin_config.encrypted_plugin", "plugins", func(value string) error {
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
