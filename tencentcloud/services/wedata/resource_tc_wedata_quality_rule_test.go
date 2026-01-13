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
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_quality_rule.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.example", "name", "tf_example_quality_rule"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.example", "type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.example", "alarm_level", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.example", "create_rule_scene", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_quality_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataQualityRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_quality_rule.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.example", "name", "tf_example_quality_rule_update"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.example", "alarm_level", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.example", "description", "Updated quality rule description"),
				),
			},
		},
	})
}

func TestAccTencentCloudNeedFixWedataQualityRuleResource_customSQL(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataQualityRuleCustomSQL,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_quality_rule.custom_sql_example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.custom_sql_example", "name", "tf_custom_sql_rule"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.custom_sql_example", "type", "3"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.custom_sql_example", "quality_dim", "1"),
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
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_quality_rule.field_config_example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.field_config_example", "name", "tf_field_config_rule"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.field_config_example", "type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule.field_config_example", "condition_type", "2"),
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
  rule_group_id                = 953
  rule_template_id             = 6142
  source_engine_types          = [2, 4, 16, 64, 128, 256, 512, 1024]
  source_object_data_type_name = "table"
  source_object_value          = "表"
  table_id                     = jsonencode(120)
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
  rule_group_id                = 953
  rule_template_id             = 6142
  source_engine_types          = [2, 4, 16, 64, 128, 256, 512, 1024]
  source_object_data_type_name = "table"
  source_object_value          = "表"
  table_id                     = jsonencode(120)
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

const testAccWedataQualityRuleCustomSQL = `
resource "tencentcloud_wedata_quality_rule" "custom_sql_example" {
  project_id        = "1948767646355341312"
  create_rule_scene = 1
  name              = "tf_custom_sql_rule"
  type              = 3
  datasource_id     = "20240101"
  database_name     = "test_db"
  table_name        = "test_table"
  quality_dim       = 1
  
  compare_rule {
    items {
      compare_type = 1
      operator     = ">"
      value_list {
        value_type = 1
        value      = "0"
      }
    }
  }
  
  alarm_level         = 2
  source_engine_types = [4]
  
  custom_sql  = "U0VMRUNUIENPVU5UKCopIEZST00gdGVzdF90YWJsZQ=="
  description = "Custom SQL quality rule"
}
`

const testAccWedataQualityRuleFieldConfig = `
resource "tencentcloud_wedata_quality_rule" "field_config_example" {
  project_id        = "1948767646355341312"
  create_rule_scene = 1
  name              = "tf_field_config_rule"
  type              = 2
  datasource_id     = "20240101"
  database_name     = "test_db"
  table_name        = "test_table"
  rule_template_id  = 2
  
  compare_rule {
    items {
      compare_type = 1
      operator     = "!="
      value_list {
        value_type = 4
        value      = "null"
      }
    }
  }
  
  alarm_level         = 1
  source_engine_types = [1]
  
  condition_type       = 2
  condition_expression = "status = 'active'"
  
  field_config {
    where_config {
      field_key       = "status"
      field_value     = "active"
      field_data_type = "string"
    }
    
    table_config {
      database_id   = "db_001"
      database_name = "test_db"
      table_id      = "tbl_001"
      table_name    = "test_table"
      table_key     = "main_table"
      
      field_config {
        field_key       = "id"
        field_value     = "primary_key"
        field_data_type = "int"
      }
    }
  }
  
  description = "Field configuration quality rule"
}
`
