package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsSyncConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_sync_config.sync_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsSyncConfig = `

resource "tencentcloud_dts_sync_config" "sync_config" {
  job_id = "sync-werwfs23"
  src_access_type = &lt;nil&gt;
  dst_access_type = &lt;nil&gt;
  options {
		init_type = "Full"
		deal_of_exist_same_table = &lt;nil&gt;
		conflict_handle_type = "ReportError"
		add_additional_column = &lt;nil&gt;
		op_types = &lt;nil&gt;
		conflict_handle_option {
			condition_column = &lt;nil&gt;
			condition_operator = &lt;nil&gt;
			condition_order_in_src_and_dst = &lt;nil&gt;
		}
		ddl_options {
			ddl_object = &lt;nil&gt;
			ddl_value = &lt;nil&gt;
		}

  }
  objects {
		mode = &lt;nil&gt;
		databases {
			db_name = &lt;nil&gt;
			new_db_name = &lt;nil&gt;
			db_mode = &lt;nil&gt;
			schema_name = &lt;nil&gt;
			new_schema_name = &lt;nil&gt;
			table_mode = &lt;nil&gt;
			tables {
				table_name = &lt;nil&gt;
				new_table_name = &lt;nil&gt;
				filter_condition = &lt;nil&gt;
			}
			view_mode = &lt;nil&gt;
			views {
				view_name = &lt;nil&gt;
				new_view_name = &lt;nil&gt;
			}
			function_mode = &lt;nil&gt;
			functions = &lt;nil&gt;
			procedure_mode = &lt;nil&gt;
			procedures = &lt;nil&gt;
			trigger_mode = &lt;nil&gt;
			triggers = &lt;nil&gt;
			event_mode = &lt;nil&gt;
			events = &lt;nil&gt;
		}
		advanced_objects = &lt;nil&gt;
		online_d_d_l = &lt;nil&gt;

  }
  job_name = &lt;nil&gt;
  job_mode = &lt;nil&gt;
  run_mode = "Immediate"
  expect_run_time = &lt;nil&gt;
  src_info {
		region = "ap-guangzhou"
		role = &lt;nil&gt;
		db_kernel = &lt;nil&gt;
		instance_id = "cdb-powiqx8q"
		ip = &lt;nil&gt;
		port = &lt;nil&gt;
		user = &lt;nil&gt;
		password = &lt;nil&gt;
		db_name = &lt;nil&gt;
		vpc_id = "vpc-92jblxto"
		subnet_id = "subnet-3paxmkdz"
		cvm_instance_id = "ins-olgl39y8"
		uniq_dcg_id = "dcg-0rxtqqxb"
		uniq_vpn_gw_id = "vpngw-9ghexg7q"
		ccn_id = "ccn-afp6kltc"
		supplier = "others"
		engine_version = "5.7"
		account = &lt;nil&gt;
		account_mode = &lt;nil&gt;
		account_role = &lt;nil&gt;
		role_external_id = &lt;nil&gt;
		tmp_secret_id = &lt;nil&gt;
		tmp_secret_key = &lt;nil&gt;
		tmp_token = &lt;nil&gt;
		encrypt_conn = "UnEncrypted"

  }
  dst_info {
		region = "ap-guangzhou"
		role = &lt;nil&gt;
		db_kernel = &lt;nil&gt;
		instance_id = "cdb-powiqx8q"
		ip = &lt;nil&gt;
		port = &lt;nil&gt;
		user = &lt;nil&gt;
		password = &lt;nil&gt;
		db_name = &lt;nil&gt;
		vpc_id = "vpc-92jblxto"
		subnet_id = "subnet-3paxmkdz"
		cvm_instance_id = "ins-olgl39y8"
		uniq_dcg_id = "dcg-0rxtqqxb"
		uniq_vpn_gw_id = "vpngw-9ghexg7q"
		ccn_id = "ccn-afp6kltc"
		supplier = "others"
		engine_version = "5.7"
		account = &lt;nil&gt;
		account_mode = &lt;nil&gt;
		account_role = &lt;nil&gt;
		role_external_id = &lt;nil&gt;
		tmp_secret_id = &lt;nil&gt;
		tmp_secret_key = &lt;nil&gt;
		tmp_token = &lt;nil&gt;
		encrypt_conn = "UnEncrypted"

  }
  auto_retry_time_range_minutes = &lt;nil&gt;
}

`
