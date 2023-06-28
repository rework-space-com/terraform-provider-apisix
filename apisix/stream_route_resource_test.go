package apisix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestStreamRouteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_upstream" "test" {
	type = "roundrobin"
	nodes = [
		{
			host   = "127.0.0.1"
			port   = 8080
			weight = 1
		}
	]
}

resource "apisix_stream_route" "test" {
	upstream_id = apisix_upstream.test.id
	remote_addr = "127.0.0.1"
	server_addr = "127.0.0.1"
	server_port = 8080
	sni         = "example.com"
}		
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// NOTE: State value checking is only necessary for Computed attributes,
					//       as the testing framework will automatically return test failures
					//       for configured attributes that mismatch the saved state.
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("apisix_stream_route.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apisix_stream_route.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "apisix_upstream" "test" {
	type = "roundrobin"
	nodes = [
		{
			host   = "127.0.0.1"
			port   = 8080
			weight = 1
		}
	]
}

resource "apisix_stream_route" "test" {
	upstream_id = apisix_upstream.test.id
	remote_addr = "127.0.0.1"
	server_addr = "127.0.0.1"
	server_port = 8080
	sni         = "example.com.ua"
}				
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apisix_stream_route.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
