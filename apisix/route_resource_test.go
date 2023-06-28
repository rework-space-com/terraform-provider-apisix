package apisix

import (
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
				// Ignore plugins value during import
				ImportStateVerifyIgnore: []string{"plugins"},
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
