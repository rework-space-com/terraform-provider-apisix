package apisix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPluginMetadataResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_plugin_metadata" "test" {
  id = "http-logger"
  metadata = jsonencode(
    {
      "log_format" = {
        "@timestamp" = "$time_iso8601"
        "client_ip"  = "$remote_addr"
        "host"       = "$host"
      }
    }
  )
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// NOTE: State value checking is only necessary for Computed attributes,
					//       as the testing framework will automatically return test failures
					//       for configured attributes that mismatch the saved state.
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("apisix_plugin_metadata.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apisix_plugin_metadata.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore metadata value during import
				ImportStateVerifyIgnore: []string{"metadata"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "apisix_plugin_metadata" "test" {
  id = "http-logger"
  metadata = jsonencode(
    {
      "log_format" = {
        "@timestamp" = "$time_iso8601"
        "client_ip"  = "$remote_addr"
      }
    }
  )
}			
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apisix_plugin_metadata.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
