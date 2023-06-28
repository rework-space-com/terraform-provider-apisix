package apisix

import (
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
