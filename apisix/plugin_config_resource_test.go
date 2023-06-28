package apisix

import (
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
				// Ignore plugins value during import
				ImportStateVerifyIgnore: []string{"plugins"},
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
