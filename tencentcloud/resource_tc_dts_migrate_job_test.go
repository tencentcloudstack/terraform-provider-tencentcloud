package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsMigrateJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateJob,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job.migrate_job", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_migrate_job.migrate_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsMigrateJob = `

resource "tencentcloud_dts_migrate_job" "migrate_job" {
  job_id = ""
  run_mode = ""
  compare_task_id = ""
  status = ""
  migrate_option {
		database_table {
			object_mode = ""
			databases {
				db_name = ""
				new_db_name = ""
				schema_name = ""
				new_schema_name = ""
				d_b_mode = ""
				schema_mode = ""
				table_mode = ""
				tables {
					table_name = &lt;nil&gt;
					new_table_name = &lt;nil&gt;
					tmp_tables = &lt;nil&gt;
					table_edit_mode = &lt;nil&gt;
				}
				view_mode = ""
				views {
					view_name = ""
					new_view_name = ""
				}
				role_mode = ""
				roles {
					role_name = ""
					new_role_name = ""
				}
				function_mode = ""
				trigger_mode = ""
				event_mode = ""
				procedure_mode = ""
				functions = 
				procedures = 
				events = 
				triggers = 
			}
			advanced_objects = 
		}
		migrate_type = ""
		consistency {
			mode = ""
		}
		is_migrate_account = 
		is_override_root = 
		is_dst_read_only = 
		extra_attr {
			key = ""
			value = ""
		}

  }
  src_info {
		region = ""
		access_type = ""
		database_type = ""
		node_type = ""
		info {
			role = ""
			db_kernel = ""
			host = ""
			port = 
			user = ""
			password = ""
			cvm_instance_id = ""
			uniq_vpn_gw_id = ""
			uniq_dcg_id = ""
			instance_id = ""
			ccn_gw_id = ""
			vpc_id = ""
			subnet_id = ""
			engine_version = ""
			account = ""
			account_role = ""
			account_mode = ""
			tmp_secret_id = ""
			tmp_secret_key = ""
			tmp_token = ""
		}
		supplier = ""
		extra_attr {
			key = ""
			value = ""
		}

  }
  dst_info {
		region = ""
		access_type = ""
		database_type = ""
		node_type = ""
		info {
			role = ""
			db_kernel = ""
			host = ""
			port = 
			user = ""
			password = ""
			cvm_instance_id = ""
			uniq_vpn_gw_id = ""
			uniq_dcg_id = ""
			instance_id = ""
			ccn_gw_id = ""
			vpc_id = ""
			subnet_id = ""
			engine_version = ""
			account = ""
			account_role = ""
			account_mode = ""
			tmp_secret_id = ""
			tmp_secret_key = ""
			tmp_token = ""
		}
		supplier = ""
		extra_attr {
			key = ""
			value = ""
		}

  }
  job_name = ""
  expect_run_time = ""
  auto_retry_time_range_minutes = 
}

`
