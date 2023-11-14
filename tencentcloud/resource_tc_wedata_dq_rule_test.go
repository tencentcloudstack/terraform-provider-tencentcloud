package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWedataDq_ruleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDq_rule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_dq_rule.dq_rule", "id")),
			},
			{
				ResourceName:      "tencentcloud_wedata_dq_rule.dq_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataDq_rule = `

resource "tencentcloud_wedata_dq_rule" "dq_rule" {
  project_id = ""
  rule_group_id = 
  name = ""
  table_id = ""
  rule_template_id = 
  type = 
  quality_dim = 
  source_object_data_type_name = ""
  source_object_value = ""
  condition_type = 
  condition_expression = ""
  custom_sql = ""
  compare_rule {
		items {
			compare_type = 
			operator = ""
			value_compute_type = 
			value_list {
				value_type = 
				value = ""
			}
		}
		cycle_step = 

  }
  alarm_level = 
  description = ""
  datasource_id = ""
  database_id = ""
  target_database_id = ""
  target_table_id = ""
  target_condition_expr = ""
  rel_condition_expr = ""
  field_config {
		where_config {
			field_key = ""
			field_value = ""
			field_data_type = ""
		}
		table_config {
			database_id = ""
			database_name = ""
			table_id = ""
			table_name = ""
			table_key = ""
			field_config {
				field_key = ""
				field_value = ""
				field_data_type = ""
			}
		}

  }
  target_object_value = ""
  source_engine_types = 
}

`
