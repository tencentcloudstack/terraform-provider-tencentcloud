package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataIntegrationTaskNodeResource_basic -v
func TestAccTencentCloudNeedFixWedataIntegrationTaskNodeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataIntegrationTaskNode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_task_node.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_task_node.example", "project_id", "1612982498218618880"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_task_node.example", "task_id", "20231022181114990"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_task_node.example", "name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_task_node.example", "node_type", "INPUT"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_task_node.example", "data_source_type", "MYSQL"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_task_node.example", "task_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_task_node.example", "task_mode"),
				),
			},
			{
				Config: testAccWedataIntegrationTaskNodeUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_task_node.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_task_node.example", "project_id", "1612982498218618880"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_task_node.example", "task_id", "20231022181114990"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_task_node.example", "name", "tf_example_update"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_task_node.example", "node_type", "INPUT"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_integration_task_node.example", "data_source_type", "MYSQL"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_task_node.example", "task_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_task_node.example", "task_mode"),
				),
			},
		},
	})
}

const testAccWedataIntegrationTaskNode = `
resource "tencentcloud_wedata_integration_task_node" "example" {
  project_id       = "1612982498218618880"
  task_id          = "20231022181114990"
  name             = "tf_example"
  node_type        = "INPUT"
  data_source_type = "MYSQL"
  task_type        = 202
  task_mode        = 2
  node_info {
    datasource_id = "5085"
    config {
      name  = "Type"
      value = "MYSQL"
    }
    config {
      name  = "splitPk"
      value = "id"
    }
    config {
      name  = "PrimaryKey"
      value = "id"
    }
    config {
      name  = "isNew"
      value = "true"
    }
    config {
      name  = "PrimaryKey_INPUT_SYMBOL"
      value = "input"
    }
    config {
      name  = "splitPk_INPUT_SYMBOL"
      value = "input"
    }
    config {
      name  = "Database"
      value = "demo_mysql"
    }
    config {
      name  = "TableNames"
      value = "users"
    }
    config {
      name  = "SiblingNodes"
      value = "[]"
    }
    schema {
      id    = "471331072"
      name  = "id"
      type  = "INT"
      alias = "id"
    }
    schema {
      id    = "422052352"
      name  = "username"
      type  = "VARCHAR(50)"
      alias = "username"
    }
  }
}
`

const testAccWedataIntegrationTaskNodeUpdate = `
resource "tencentcloud_wedata_integration_task_node" "example" {
  project_id       = "1612982498218618880"
  task_id          = "20231022181114990"
  name             = "tf_example_update"
  node_type        = "INPUT"
  data_source_type = "MYSQL"
  task_type        = 202
  task_mode        = 2
  node_info {
    datasource_id = "5085"
    config {
      name  = "Type"
      value = "MYSQL"
    }
    config {
      name  = "splitPk"
      value = "id"
    }
    config {
      name  = "PrimaryKey"
      value = "id"
    }
    config {
      name  = "isNew"
      value = "true"
    }
    config {
      name  = "PrimaryKey_INPUT_SYMBOL"
      value = "input"
    }
    config {
      name  = "splitPk_INPUT_SYMBOL"
      value = "input"
    }
    config {
      name  = "Database"
      value = "demo_mysql"
    }
    config {
      name  = "TableNames"
      value = "users"
    }
    config {
      name  = "SiblingNodes"
      value = "[]"
    }
    schema {
      id    = "471331072"
      name  = "id"
      type  = "INT"
      alias = "id"
    }
    schema {
      id    = "422052352"
      name  = "username"
      type  = "VARCHAR(50)"
      alias = "username"
    }
  }
}
`
