package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_access_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "dst_access_type"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "job_name", "tf_test_sync_config"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "job_mode", "liteMode"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "run_mode", "Immediate"),
					// objects
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.mode", "Partial"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.db_name", "tf_ci_test"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.new_db_name", "tf_ci_test_new"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.db_mode", "Partial"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.table_mode", "All"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.tables.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.tables.0.table_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "objects.0.databases.0.tables.0.new_table_name", "test_new"),
					// src_info dest_info
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.0.instance_id", "cdb-fitq5t9h"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.0.user", "keep_dts"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "src_info.0.db_name", "tf_ci_test"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_info.0.vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "src_info.0.subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "dst_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "dst_info.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "dst_info.0.instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "dst_info.0.db_name", "tf_ci_test_new"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "dst_info.0.vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_config.sync_config", "dst_info.0.subnet_id"),

					resource.TestCheckResourceAttr("tencentcloud_dts_sync_config.sync_config", "auto_retry_time_range_minutes", "0"),
				),
			},
			{
				ResourceName:            "tencentcloud_dts_sync_config.sync_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dst_info.0.password", "src_info.0.password", "job_mode"},
			},
		},
	})
}

const testAccDtsSyncConfig = testAccDtsMigrateJob_vpc_config + `
resource "tencentcloud_cynosdb_cluster" "foo" {
	available_zone               = var.availability_zone
	vpc_id                       = local.vpc_id
	subnet_id                    = local.subnet_id
	db_type                      = "MYSQL"
	db_version                   = "5.7"
	storage_limit                = 1000
	cluster_name                 = "tf-cynosdb-mysql-sync-dst"
	password                     = "cynos@123"
	instance_maintain_duration   = 3600
	instance_maintain_start_time = 10800
	instance_maintain_weekdays   = [
	  "Fri",
	  "Mon",
	  "Sat",
	  "Sun",
	  "Thu",
	  "Wed",
	  "Tue",
	]
  
	instance_cpu_core    = 1
	instance_memory_size = 2
	param_items {
	  name = "character_set_server"
	  current_value = "utf8"
	}
	param_items {
	  name = "time_zone"
	  current_value = "+09:00"
	}
	param_items {
		name = "lower_case_table_names"
		current_value = "1"
	}
  
	force_delete = true
  
	rw_group_sg = [
	  local.sg_id
	]
	ro_group_sg = [
	  local.sg_id
	]
	prarm_template_id = var.my_param_template
  }

resource "tencentcloud_dts_sync_job" "sync_job" {
	pay_mode = "PostPay"
	src_database_type = "mysql"
	src_region = "ap-guangzhou"
	dst_database_type = "cynosdbmysql"
	dst_region = "ap-guangzhou"
	tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
	auto_renew = 0
	instance_class = "micro"
  }

resource "tencentcloud_dts_sync_config" "sync_config" {
  job_id = tencentcloud_dts_sync_job.sync_job.job_id
  src_access_type = "cdb"
  dst_access_type = "cdb"
  
  job_name = "tf_test_sync_config"
  job_mode = "liteMode"
  run_mode = "Immediate"

  objects {
	mode = "Partial"
      databases {
	    db_name = "tf_ci_test"
			new_db_name = "tf_ci_test_new"
			db_mode = "Partial"
			table_mode = "All"
			tables {
				table_name = "test"
				new_table_name = "test_new"
			}
	  }
  }
  src_info {
		region        = "ap-guangzhou"
		instance_id   = "cdb-fitq5t9h"
		user          = "keep_dts"
		password      = "Letmein123"
		db_name       = "tf_ci_test"
		vpc_id        = local.vpc_id
		subnet_id     = local.subnet_id
  }
  dst_info {
		region        = "ap-guangzhou"
		instance_id   = tencentcloud_cynosdb_cluster.foo.id
		user          = "root"
		password      = "cynos@123"
		db_name       = "tf_ci_test_new"
		vpc_id        = local.vpc_id
		subnet_id     = local.subnet_id
  }
  auto_retry_time_range_minutes = 0
}

`
