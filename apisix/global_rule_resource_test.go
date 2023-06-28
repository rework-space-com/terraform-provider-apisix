package apisix

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGlobalRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "apisix_global_rule" "test" {
	id = "123"
	plugins = jsonencode(
		{
			prometheus = {
				prefer_name = true
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
					resource.TestCheckResourceAttrSet("apisix_global_rule.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apisix_global_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore plugins value during import
				ImportStateVerifyIgnore: []string{"plugins"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "apisix_global_rule" "test" {
	id = "123"
	plugins = jsonencode(
		{
			prometheus = {
				prefer_name = false
			}
		}
	)
}			
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apisix_global_rule.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
