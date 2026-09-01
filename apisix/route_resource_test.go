package apisix

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRouteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_route" "test" {
	name         = "Example"
	desc         = "Example of the route configuration"
	uris         = ["/api/v1", "/status"]
	hosts        = ["foo.com", "*.bar.com"]
	remote_addrs = ["10.0.0.0/8"]
	methods      = ["GET", "POST"]
	vars = jsonencode(
		[["http_user", "==", "ios"]]
	)
	timeout = {
		connect = 3
		send    = 3
		read    = 3
	}
	plugins = jsonencode(
		{
			ip-restriction = {
				blacklist = ["10.10.10.0/24"]
				message   = "Access denied"
			}
		}
	)
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
					resource.TestCheckResourceAttrSet("apisix_route.test", "id"),
					resource.TestCheckResourceAttr("apisix_route.test", "priority", "0"),
					resource.TestCheckResourceAttr("apisix_route.test", "priority", "0"),
					resource.TestCheckResourceAttr("apisix_route.test", "enable_websocket", "false"),
					resource.TestCheckResourceAttr("apisix_route.test", "status", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apisix_route.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "apisix_route" "test" {
					name         = "Example"
					desc         = "Example of the route configuration"
					uris         = ["/api/v1", "/status"]
					hosts        = ["foo.com"]
					remote_addrs = ["10.0.0.0/8"]
					methods      = ["GET", "POST", "PUT"]
					vars = jsonencode(
						[["http_user", "==", "ios"]]
					)
					timeout = {
						connect = 10
						send    = 5
						read    = 10
					}
					plugins = jsonencode(
						{
							ip-restriction = {
								blacklist = ["10.10.10.0/24"]
								message   = "Access denied"
							}
						}
					)
					labels = {
						"version" = "v2"
					}
					enable_websocket = true
				}				
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apisix_route.test", "id"),
					resource.TestCheckResourceAttr("apisix_route.test", "priority", "0"),
					resource.TestCheckResourceAttr("apisix_route.test", "priority", "0"),
					resource.TestCheckResourceAttr("apisix_route.test", "enable_websocket", "true"),
					resource.TestCheckResourceAttr("apisix_route.test", "status", "1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestRouteResource_EncryptedPluginField is a regression test for #41: APISIX
// returns encrypt_fields (e.g. openid-connect's client_secret) as ciphertext
// from the Create response, but as plaintext on a subsequent GET. Create()
// must re-fetch in that case, or the reported new state won't match the plan
// and Terraform fails with "Provider produced inconsistent result after
// apply", tainting the resource and looping on every following apply.
func TestRouteResource_EncryptedPluginField(t *testing.T) {
	config := providerConfig + `
resource "apisix_route" "encrypted_plugin" {
	name  = "EncryptedPluginField"
	uris  = ["/encrypted-plugin-field"]
	hosts = ["encrypted-plugin-field.example.com"]
	plugins = jsonencode(
		{
			openid-connect = {
				bearer_only   = true
				unauth_action = "pass"
				client_id     = "example-client-id"
				client_secret = "PLAINTEXT-SECRET-VALUE-12345"
				discovery     = "https://example.com/.well-known/openid-configuration"
			}
		}
	)
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apisix_route.encrypted_plugin", "id"),
					// The actual regression check: state must hold the plaintext
					// plugins config we configured, not the ciphertext APISIX's
					// Create response returns for encrypt_fields. Before the fix,
					// this would already have failed the apply itself with
					// "Provider produced inconsistent result after apply" -- this
					// assertion pins down what the correct state value must be.
					resource.TestCheckResourceAttrWith("apisix_route.encrypted_plugin", "plugins", func(value string) error {
						var actual interface{}
						if err := json.Unmarshal([]byte(value), &actual); err != nil {
							return fmt.Errorf("failed to unmarshal plugins: %w", err)
						}

						expected := map[string]interface{}{
							"openid-connect": map[string]interface{}{
								"bearer_only":   true,
								"unauth_action": "pass",
								"client_id":     "example-client-id",
								"client_secret": "PLAINTEXT-SECRET-VALUE-12345",
								"discovery":     "https://example.com/.well-known/openid-configuration",
							},
						}

						if !reflect.DeepEqual(actual, expected) {
							return fmt.Errorf("expected plugins to equal %#v, got %#v (state likely holds the ciphertext APISIX's Create response returned instead of the plaintext client_secret)", expected, actual)
						}
						return nil
					}),
				),
			},
			// Plan-only with the identical config must show no changes. Before
			// the fix, the resource was tainted right after apply, so this step
			// would have proposed a destroy+recreate.
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestRouteResource_EncryptedPluginFieldUpdate covers the Update path with an
// encrypt_field (openid-connect client_secret): an in-place Update re-PUTs the
// identical field. Update() must set state.plugins to the planned value, or Terraform
// fails with "Provider produced inconsistent result after apply".
// NOTE: same stable-keyring caveat as TestServiceResource_EncryptedPluginField.
func TestRouteResource_EncryptedPluginFieldUpdate(t *testing.T) {
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
resource "apisix_route" "encrypted_plugin_update" {
	name  = "EncryptedPluginFieldUpdate"
	uris  = ["/encrypted-plugin-field-update"]
	hosts = ["encrypted-plugin-field-update.example.com"]
	plugins = jsonencode(
		{
` + pluginsBlock + `
		}
	)
}
`

	// In-place Update (enable_websocket) re-PUTs the identical encrypted plugin field.
	updatedConfig := providerConfig + `
resource "apisix_route" "encrypted_plugin_update" {
	name  = "EncryptedPluginFieldUpdate"
	uris  = ["/encrypted-plugin-field-update"]
	hosts = ["encrypted-plugin-field-update.example.com"]
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
		resource.TestCheckResourceAttrSet("apisix_route.encrypted_plugin_update", "id"),
		resource.TestCheckResourceAttrWith("apisix_route.encrypted_plugin_update", "plugins", func(value string) error {
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
