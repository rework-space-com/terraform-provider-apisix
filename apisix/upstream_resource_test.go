package apisix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUpstreamResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_upstream" "test" {
	name = "Example"
	desc = "Example of the upstream resource usage"
	type = "roundrobin"
	labels = {
		version = "v1"
	}
	nodes = [
		{
			host   = "127.0.0.1"
			port   = 1980
			weight = 1
		},
		{
			host   = "127.0.0.1"
			port   = 1970
			weight = 1
		},
	]
	keepalive_pool = {
		idle_timeout = 5
		requests     = 10
		size         = 15
	}
	checks = {
		active = {
			host      = "example.com"
			port      = 8888
			timeout   = 5
			http_path = "/status"
			healthy = {
				interval  = 2,
				successes = 1
			}
			unhealthy = {
				interval      = 1
				http_failures = 2
			}
		}
		passive = {
			healthy = {
				http_statuses = [200, 201]
			}
			unhealthy = {
				http_statuses = [500]
				http_failures = 3
				tcp_failures  = 3
			}
		}
	}
}		
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// NOTE: State value checking is only necessary for Computed attributes,
					//       as the testing framework will automatically return test failures
					//       for configured attributes that mismatch the saved state.
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("apisix_upstream.test", "id"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "type", "roundrobin"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "pass_host", "pass"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "scheme", "http"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "hash_on", "vars"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apisix_upstream.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "apisix_upstream" "test" {
	name = "Example"
	desc = "Example of the upstream resource usage"
	type = "roundrobin"
	labels = {
		version = "v2"
	}
	nodes = [
		{
			host   = "127.0.0.1"
			port   = 1980
			weight = 1
		},
		{
			host   = "127.0.0.1"
			port   = 1970
			weight = 5
		},
	]
	keepalive_pool = {
		idle_timeout = 10
		requests     = 10
		size         = 15
	}
	checks = {
		active = {
			host      = "example.com"
			port      = 8888
			timeout   = 5
			http_path = "/status"
			healthy = {
				interval  = 3,
				successes = 1
			}
			unhealthy = {
				interval      = 1
				http_failures = 2
			}
		}
	}
}				
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apisix_upstream.test", "id"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "type", "roundrobin"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "pass_host", "pass"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "scheme", "http"),
					resource.TestCheckResourceAttr("apisix_upstream.test", "hash_on", "vars"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
