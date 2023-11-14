package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWedataDq_ruleTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDq_ruleTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.dq_rule_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_wedata_dq_rule_template.dq_rule_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataDq_ruleTemplate = `

resource "tencentcloud_wedata_dq_rule_template" "dq_rule_template" {
  type = 
  name = ""
  quality_dim = 
  source_object_type = 
  description = ""
  source_engine_types = 
  multi_source_flag = 
  sql_expression = ""
  project_id = ""
  where_flag = 
}

`
