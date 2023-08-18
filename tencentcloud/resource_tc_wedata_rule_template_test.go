package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudWedataRuleTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataRuleTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_rule_template.rule_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_wedata_rule_template.rule_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataRuleTemplate = `

resource "tencentcloud_wedata_rule_template" "rule_template" {
  type                = 2
  name                = "fo test"
  quality_dim         = 3
  source_object_type  = 2
  description         = "for tf test"
  source_engine_types = [3]
  multi_source_flag   = false
  sql_expression      = "c2VsZWN0ICogZnJvbSBkYg=="
  project_id          = "1840731346428280832"
  where_flag          = false
}

`
