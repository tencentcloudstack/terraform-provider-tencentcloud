package wedata_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataDqRuleResource_basic -v
func TestAccTencentCloudNeedFixWedataDqRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDqRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_dq_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataDqRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule.example", "id"),
				),
			},
		},
	})
}

const testAccWedataDqRule = `
resource "tencentcloud_wedata_dq_rule" "example" {
  project_id                   = "1948767646355341312"
  rule_group_id                = 312
  rule_template_id             = 1
  name                         = "tf_example"
  table_id                     = "N85hbsh5QQ2VLHL2iOUVeQ"
  type                         = 1
  source_object_data_type_name = "table"
  source_object_value          = "表"
  condition_type               = 1
  compare_rule {
    items {
      compare_type = 1
      operator     = "=="
      value_list {
        value_type = 3
        value      = "100"
      }
    }
  }
  alarm_level = 1
  description = "description."
}
`

const testAccWedataDqRuleUpdate = `
resource "tencentcloud_wedata_dq_rule" "example" {
  project_id                   = "1948767646355341312"
  rule_group_id                = 312
  rule_template_id             = 16
  name                         = "tf_example_update"
  table_id                     = "N85hbsh5QQ2VLHL2iOUVeQ"
  type                         = 1
  condition_expression         = "${key1}"
  source_object_data_type_name = "string"
  source_object_value          = "id"
  condition_type               = 2
  compare_rule {
    items {
      compare_type       = 1
      value_compute_type = 1
      value_list {
        value_type = 1
        value      = "50"
      }
      value_list {
        value_type = 2
        value      = "80"
      }
    }
  }
  alarm_level = 2
  description = "description update."
}
`
