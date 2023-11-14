package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWedataIntegration_real_time_taskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataIntegration_real_time_task,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_real_time_task.integration_real_time_task", "id")),
			},
			{
				ResourceName:      "tencentcloud_wedata_integration_real_time_task.integration_real_time_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataIntegration_real_time_task = `

resource "tencentcloud_wedata_integration_real_time_task" "integration_real_time_task" {
  task_info {
		task_name = "TaskTest_10"
		description = "Task for test"
		sync_type = 2
		task_type = 201
		workflow_id = "1"
		task_id = "j84cc717e-215b-4960-9575-898586bae37f"
		schedule_task_id = "1"
		task_group_id = "1"
		project_id = "1455251608631480391"
		creator_uin = "100028448000"
		operator_uin = "100028448000"
		owner_uin = "100028448000"
		app_id = "1315000000"
		status = 1
		nodes {
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
		executor_id = "2000"
		config {
			name = "Database"
			value = "db"
		}
		ext_config {
			name = "x"
			value = "320"
		}
		execute_context {
			name = ""
			value = ""
		}
		mappings {
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
		task_mode = "1"
		incharge = ""
		offline_task_add_entity {
			workflow_name = "workflow_test"
			dependency_workflow = "no"
			start_time = "2023-12-31 00:00:00"
			end_time = "2099-12-31 00:00:00"
			cycle_type = 0
			cycle_step = 1
			delay_time = 0
			crontab_expression = "0 0 1 * * ?"
			retry_wait = 5
			retriable = 1
			try_limit = 1
			run_priority = 6
			product_name = "DATA_INTEGRATION"
			self_depend = 3
			task_action = "1"
			execution_end_time = "16:59"
			execution_start_time = "02:00"
			task_auto_submit = false
			instance_init_strategy = &lt;nil&gt;
		}
		executor_group_name = "executor1"
		in_long_manager_url = "172.16.0.3:8083"
		in_long_stream_id = "b_q3b502073-1cac-4a7b-a67f-d30314833a32"
		in_long_manager_version = "v16"
		data_proxy_url = 
		submit = false
		input_datasource_type = "MYSQL"
		output_datasource_type = "MYSQL"
		num_records_in = 1000
		num_records_out = 1000
		reader_delay = 
		num_restarts = 1
		create_time = "2023-10-12 17:17:14"
		update_time = "2023-10-12 17:17:14"
		last_run_time = "2023-10-12 17:17:14"
		stop_time = "2023-10-12 17:17:14"
		has_version = false
		locked = false
		locker = "100028578868"
		running_cu = 
		task_alarm_regular_list = &lt;nil&gt;
		switch_resource = 0
		read_phase = 0
		instance_version = 1

  }
  project_id = "1455251608631480391"
}

`
