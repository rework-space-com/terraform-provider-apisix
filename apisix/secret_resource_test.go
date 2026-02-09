package apisix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSecretResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_secret" "test" {
	id = "123"
	gcp = {
		auth_file  = "/path/to/file.json"
		ssl_verify = false
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// NOTE: State value checking is only necessary for Computed attributes,
					//       as the testing framework will automatically return test failures
					//       for configured attributes that mismatch the saved state.
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("apisix_secret.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apisix_secret.test",
				ImportState:       true,
				ImportStateId:     "gcp/123",
				ImportStateVerify: true,
				// Ignore plugins value during import
				ImportStateVerifyIgnore: []string{"plugins"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "apisix_secret" "test" {
	id = "123"
	gcp = {
		auth_file  = "/path/to/file.json"
		ssl_verify = false
	}
}			
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apisix_secret.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
