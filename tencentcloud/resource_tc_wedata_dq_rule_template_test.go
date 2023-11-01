package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataDqRuleTemplateResource_basic -v
func TestAccTencentCloudNeedFixWedataDqRuleTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDqRuleTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "quality_dim"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "source_object_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "multi_source_flag"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "sql_expression"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "where_flag"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_dq_rule_template.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataDqRuleTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "quality_dim"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "source_object_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "multi_source_flag"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "sql_expression"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule_template.example", "where_flag"),
				),
			},
		},
	})
}

const testAccWedataDqRuleTemplate = `
resource "tencentcloud_wedata_dq_rule_template" "example" {
  type                = 2
  name                = "tf_example"
  quality_dim         = 1
  source_object_type  = 2
  description         = "description."
  source_engine_types = [2]
  multi_source_flag   = true
  sql_expression      = "c2VsZWN0"
  project_id          = "1948767646355341312"
  where_flag          = true
}
`

const testAccWedataDqRuleTemplateUpdate = `
resource "tencentcloud_wedata_dq_rule_template" "example" {
  type                = 2
  name                = "tf_example"
  quality_dim         = 1
  source_object_type  = 2
  description         = "description."
  source_engine_types = [2]
  multi_source_flag   = true
  sql_expression      = "c2VsZWN0"
  project_id          = "1948767646355341312"
  where_flag          = true
}
`
