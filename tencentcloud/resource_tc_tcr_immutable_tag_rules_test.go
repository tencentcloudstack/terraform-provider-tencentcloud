package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrImmutableTagRulesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrImmutableTagRules,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_immutable_tag_rules.immutable_tag_rules", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_immutable_tag_rules.immutable_tag_rules",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrImmutableTagRules = `

resource "tencentcloud_tcr_immutable_tag_rules" "immutable_tag_rules" {
  registry_id = "tcr-xxx"
  namespace_name = "ns"
  rule {
		repository_pattern = "**"
		tag_pattern = "**"
		repository_decoration = "repoMatches"
		tag_decoration = "matches"
		disabled = false
		rule_id = 1
		ns_name = "ns"

  }
  tags = {
    "createdBy" = "terraform"
  }
}

`
