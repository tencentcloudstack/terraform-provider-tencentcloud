package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWedataIntegration_task_nodeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataIntegration_task_node,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_task_node.integration_task_node", "id")),
			},
			{
				ResourceName:      "tencentcloud_wedata_integration_task_node.integration_task_node",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataIntegration_task_node = `

resource "tencentcloud_wedata_integration_task_node" "integration_task_node" {
  node_info {
		id = ""
		task_id = "j84cc717e-215b-4960-9575-898586bae37f"
		name = "input_name"
		node_type = "INPUT"
		data_source_type = "MYSQL"
		description = "Node for test"
		datasource_id = "100"
		config {
			name = "Database"
			value = "db"
		}
		ext_config {
			name = "x"
			value = "320"
		}
		schema {
			id = "796598528"
			name = "col_name"
			type = "string"
			value = "1"
			properties {
				name = "name"
				value = "value"
			}
			alias = "name"
			comment = "comment"
		}
		node_mapping {
			source_id = "10"
			sink_id = "11"
			source_schema {
				id = "796598528"
				name = "col_name"
				type = "string"
				value = "1"
				properties {
					name = "name"
					value = "value"
				}
				alias = "name"
				comment = "comment"
			}
			schema_mappings {
				source_schema_id = "200"
				sink_schema_id = "300"
			}
			ext_config {
				name = "x"
				value = "320"
			}
		}
		app_id = "1315000000"
		project_id = "1455251608631480391"
		creator_uin = "100028448000"
		operator_uin = "100028448000"
		owner_uin = "100028448000"
		create_time = "2023-10-17 18:02:46"
		update_time = "2023-10-17 18:02:46"

  }
  project_id = "1455251608631480391"
  task_type = 201
}

`
