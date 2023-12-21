package dts_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDtsSyncJobStartOperationResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncJobOperation_start("operation_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_start_operation.sync_job_start_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_start_operation.sync_job_start_operation", "job_id"),
				),
			},
			{
				Config: testAccDtsSyncJobOperation_pause("operation_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_pause_operation.sync_job_pause_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_pause_operation.sync_job_pause_operation", "job_id"),
				),
			},
			{
				Config: testAccDtsSyncJobOperation_continue("operation_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_continue_operation.sync_job_continue_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_continue_operation.sync_job_continue_operation", "job_id"),
				),
			},
			{
				Config: testAccDtsSyncJobOperation_isolate("operation_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_isolate_operation.sync_job_isolate_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_isolate_operation.sync_job_isolate_operation", "job_id"),
				),
			},
			{
				Config: testAccDtsSyncJobOperation_recover("operation_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_recover_operation.sync_job_recover_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_recover_operation.sync_job_recover_operation", "job_id"),
				),
			},
			{
				Config: testAccDtsSyncJobOperation_stop("operation_basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_stop_operation.sync_job_stop_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_stop_operation.sync_job_stop_operation", "job_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudDtsSyncJobStartOperationResource_resize(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncJobOperation_start("operation_resize"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_start_operation.sync_job_start_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_start_operation.sync_job_start_operation", "job_id"),
				),
			},
			{
				Config: testAccDtsSyncJobOperation_resize("operation_resize"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_resize_operation.sync_job_resize_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_resize_operation.sync_job_resize_operation", "job_id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_job_resize_operation.sync_job_resize_operation", "new_instance_class", "micro"),
				),
			},
		},
	})
}

func testAccDtsSyncJobOperation_basic(name string) string {
	ret := fmt.Sprintf(testAccDtsMigrateJob_vpc_config+`
	resource "tencentcloud_cynosdb_cluster" "foo" {
		available_zone               = var.availability_zone
		vpc_id                       = local.vpc_id
		subnet_id                    = local.subnet_id
		db_type                      = "MYSQL"
		db_version                   = "5.7"
		storage_limit                = 1000
		cluster_name                 = "tf-cynosdb-dts-sync-job-src-%s"
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
	  
	  job_name = "tf_test_sync_config_%s"
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

	resource "tencentcloud_dts_sync_check_job_operation" "sync_check_job_operation" {
		job_id = tencentcloud_dts_sync_config.sync_config.job_id
	}

	locals {
		job_id_checked = tencentcloud_dts_sync_check_job_operation.sync_check_job_operation.id
	}
	
`, name, name)
	return ret
}

func testAccDtsSyncJobOperation_start(name string) string {

	ret := testAccDtsSyncJobOperation_basic(name) + `

	resource "tencentcloud_dts_sync_job_start_operation" "sync_job_start_operation" {
		job_id = local.job_id_checked
	}
`
	return ret
}

func testAccDtsSyncJobOperation_stop(name string) string {

	ret := testAccDtsSyncJobOperation_basic(name) + `

	resource "tencentcloud_dts_sync_job_stop_operation" "sync_job_stop_operation" {
		job_id = local.job_id_checked
	}
`
	return ret
}

func testAccDtsSyncJobOperation_pause(name string) string {

	ret := testAccDtsSyncJobOperation_basic(name) + `

	resource "tencentcloud_dts_sync_job_pause_operation" "sync_job_pause_operation" {
    	job_id = local.job_id_checked
	}
`
	return ret
}

// func testAccDtsSyncJobOperation_resume(name string) string {

// 	ret := testAccDtsSyncJobOperation_basic(name) + `

// 	resource "tencentcloud_dts_sync_job_resume_operation" "sync_job_resume_operation" {
// 		job_id = local.job_id_checked
// 	}
// `
// 	return ret
// }

func testAccDtsSyncJobOperation_resize(name string) string {

	ret := testAccDtsSyncJobOperation_basic(name) + `

	resource "tencentcloud_dts_sync_job_resize_operation" "sync_job_resize_operation" {
    	job_id = local.job_id_checked
		new_instance_class = "micro"
	}
`
	return ret
}

func testAccDtsSyncJobOperation_recover(name string) string {

	ret := testAccDtsSyncJobOperation_basic(name) + `

	resource "tencentcloud_dts_sync_job_recover_operation" "sync_job_recover_operation" {
    	job_id = local.job_id_checked
	}
`
	return ret
}

func testAccDtsSyncJobOperation_isolate(name string) string {

	ret := testAccDtsSyncJobOperation_basic(name) + `

	resource "tencentcloud_dts_sync_job_isolate_operation" "sync_job_isolate_operation" {
    	job_id = local.job_id_checked
	}
`
	return ret
}

func testAccDtsSyncJobOperation_continue(name string) string {

	ret := testAccDtsSyncJobOperation_basic(name) + `
	resource "tencentcloud_dts_sync_job_continue_operation" "sync_job_continue_operation" {
    	job_id = local.job_id_checked
	}
`
	return ret
}
