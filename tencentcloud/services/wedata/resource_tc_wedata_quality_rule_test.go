package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataQualityRuleResource_basic -v
func TestAccTencentCloudNeedFixWedataQualityRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataQualityRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_quality_rule.system_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.system_rule", "name", "system_template_rule_tf_test"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.system_rule", "type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.system_rule", "alarm_level", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.system_rule", "create_rule_scene", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.system_rule", "database_name", "default"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.system_rule", "table_name", "big_table_500"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_quality_rule.system_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataQualityRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_quality_rule.system_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.system_rule", "name", "system_template_rule_tf_test"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.system_rule", "alarm_level", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.system_rule", "description", "rule desc"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.system_rule", "condition_expression", "dt='000'"),
				),
			},
		},
	})
}

func TestAccTencentCloudNeedFixWedataQualityRuleResource_fieldConfig(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataQualityRuleFieldConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_quality_rule.custom_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.custom_rule", "name", "custom_template_rule_tf_test"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.custom_rule", "type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.custom_rule", "condition_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.custom_rule", "alarm_level", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.custom_rule", "database_name", "default"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.custom_rule", "table_name", "at_src_mysql2hive_prod_cq_makeup_09db_1_di"),
				),
			},
		},
	})
}

const testAccWedataQualityRule = `
resource "tencentcloud_wedata_quality_rule" "system_rule" {
  alarm_level                  = 2
  condition_expression         = jsonencode(222)
  condition_type               = 2
  create_rule_scene            = 1
  database_name                = "default"
  datasource_id                = jsonencode(65253)
  description                  = jsonencode(111)
  name                         = "system_template_rule_tf_test"
  project_id                   = jsonencode(3016337760439783424)
  quality_dim                  = 1
  # rule_group_id                = 953
  rule_template_id             = 6142
  source_engine_types          = [2, 4, 16, 64, 128, 256, 512, 1024]
  source_object_data_type_name = "table"
  source_object_value          = "表"
  table_id                     = 176
  table_name                   = "big_table_500"
  type                         = 1
  compare_rule {
    compute_expression = jsonencode(0)
    cycle_step         = 0
    items {
      compare_type       = 1
      operator           = "<"
      value_compute_type = 0
      value_list {
        value      = jsonencode(100)
        value_type = 3
      }
    }
  }
}

`

const testAccWedataQualityRuleUpdate = `
resource "tencentcloud_wedata_quality_rule" "system_rule" {
  alarm_level                  = 1
  condition_expression         = "dt='000'"
  condition_type               = 2
  create_rule_scene            = 1
  database_name                = "default"
  datasource_id                = 65253
  description                  = "rule desc"
  name                         = "system_template_rule_tf_test"
  project_id                   = jsonencode(3016337760439783424)
  quality_dim                  = 1
  # rule_group_id                = 953
  rule_template_id             = 6142
  source_engine_types          = [2, 4, 16, 64, 128, 256, 512, 1024]
  source_object_data_type_name = "table"
  source_object_value          = "表"
  table_id                     = 176
  table_name                   = "big_table_500"
  type                         = 1
  compare_rule {
    compute_expression = jsonencode(0)
    cycle_step         = 0
    items {
      compare_type       = 1
      operator           = "<="
      value_compute_type = 0
      value_list {
        value      = jsonencode(100)
        value_type = 3
      }
    }
  }
}
`

const testAccWedataQualityRuleFieldConfig = `
resource "tencentcloud_wedata_quality_rule" "custom_rule" {
  alarm_level                  = 2
  condition_type               = 1
  create_rule_scene            = 1
  database_name                = "default"
  datasource_id                = jsonencode(65253)
  description                  = jsonencode(111)
  index                        = null
  name                         = "custom_template_rule_tf_test"
  project_id                   = jsonencode(3016337760439783424)
  quality_dim                  = 1
  # rule_group_id                = 949
  rule_template_id             = 6809
  source_engine_types          = [2, 4, 16]
  source_object_data_type_name = "table"
  source_object_value          = "表"
  table_id                     = jsonencode(175)
  table_name                   = "at_src_mysql2hive_prod_cq_makeup_09db_1_di"
  type                         = 2
  compare_rule {
    compute_expression = jsonencode(0)
    cycle_step         = 0
    items {
      compare_type       = 1
      operator           = "<"
      value_compute_type = 0
      value_list {
        value      = jsonencode(100)
        value_type = 3
      }
    }
  }
  field_config {
    table_config {
      database_id   = null
      database_name = "default"
      table_id      = null
      table_key     = "table_1"
      table_name    = "at_src_mysql2hive_prod_cq_makeup_09db_1_di"
    }
  }
}
`
