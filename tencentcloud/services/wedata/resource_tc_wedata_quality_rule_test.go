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
resource "tencentcloud_wedata_quality_rule" "example" {
  project_id        = "1948767646355341312"
  create_rule_scene = 1
  name              = "tf_example_quality_rule"
  type              = 1
  datasource_id     = "20240101"
  database_name     = "test_db"
  table_name        = "test_table"
  rule_template_id  = 1
  quality_dim       = 1
  rule_group_id     = 312
  
  compare_rule {
    items {
      compare_type = 1
      operator     = "=="
      value_list {
        value_type = 3
        value      = "100"
      }
    }
    cycle_step         = 3600
    compute_expression = "1"
  }
  
  alarm_level         = 1
  source_engine_types = [1, 2]
  
  source_object_data_type_name = "table"
  source_object_value          = "test_table"
  condition_type               = 1
  
  description = "Example quality rule for testing"
}
`

const testAccWedataQualityRuleUpdate = `
resource "tencentcloud_wedata_quality_rule" "example" {
  project_id        = "1948767646355341312"
  create_rule_scene = 1
  name              = "tf_example_quality_rule_update"
  type              = 1
  datasource_id     = "20240101"
  database_name     = "test_db"
  table_name        = "test_table"
  rule_template_id  = 16
  quality_dim       = 2
  rule_group_id     = 312
  
  compare_rule {
    items {
      compare_type       = 2
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
    cycle_step         = 7200
    compute_expression = "1"
  }
  
  alarm_level         = 2
  source_engine_types = [1, 4]
  
  source_object_data_type_name = "string"
  source_object_value          = "id"
  condition_type               = 2
  condition_expression         = "status = 'active'"
  
  description = "Updated quality rule description"
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
