package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudDtsMigrateJobResource_basic(t *testing.T) {
	t.Parallel()
	curSec := fmt.Sprint(time.Now().Unix())
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDtsMigrateJobDestroy,
		Steps: []resource.TestStep{
			{
				// PreventDiskCleanup: true,
				Config: fmt.Sprintf(testAccDtsMigrateJob_basic, "migrate_service_1", curSec),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDtsMigrateJobExists("tencentcloud_dts_migrate_job.job"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job.job", "service_id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job.job", "run_mode", "immediate"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job.job", "migrate_option.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job.job", "src_info.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job.job", "dst_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job.job", "auto_retry_time_range_minutes", "0"),
				),
			},
			{
				ResourceName:            "tencentcloud_dts_migrate_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_retry_time_range_minutes", "dst_info", "src_info"}, // dts don't support to query auto_retry_time_range_minutes in Read
			},
		},
	})
}

func testAccCheckDtsMigrateJobDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	dtsService := DtsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dts_migrate_job" {
			continue
		}

		job, err := dtsService.DescribeDtsMigrateJobById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if job != nil {
			status := *job.TradeInfo.TradeStatus
			if status != "isolated" && status != "offlined" {
				return fmt.Errorf("DTS migrate job still exist, Id: %v, status:%s", rs.Primary.ID, status)
			}
		}
	}
	return nil
}

func testAccCheckDtsMigrateJobExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		dtsService := DtsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("DTS migrate job %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("DTS migrate job id is not set")
		}

		job, err := dtsService.DescribeDtsMigrateJobById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if job == nil {
			return fmt.Errorf("DTS migrate job not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccDtsMigrateJob_vpc_config = `
data "tencentcloud_security_groups" "internal" {
	name = "default"
}

data "tencentcloud_vpc_instances" "vpc" {
	name ="Default-VPC"
}
	
data "tencentcloud_vpc_subnets" "subnet" {
	vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
}
	
locals {
	vpc_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.vpc_id
	subnet_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id
	sg_id = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
}

variable "availability_zone" {
	default = "ap-guangzhou-4"
  }
  
variable "my_param_template" {
	  default = "15765"
  }
`

const testAccDtsMigrateJob_cynosdb_mysql = testAccDtsMigrateJob_vpc_config + `
resource "tencentcloud_cynosdb_cluster" "foo" {
	available_zone               = var.availability_zone
	vpc_id                       = local.vpc_id
	subnet_id                    = local.subnet_id
	db_type                      = "MYSQL"
	db_version                   = "5.7"
	storage_limit                = 1000
	cluster_name                 = "tf-cynosdb-mysql"
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
`

const testAccDtsMigrateJob_basic = testAccDtsMigrateJob_cynosdb_mysql + `

resource "tencentcloud_dts_migrate_service" "service" {
	src_database_type = "mysql"
	dst_database_type = "cynosdbmysql"
	src_region = "ap-guangzhou"
	dst_region = "ap-guangzhou"
	instance_class = "small"
	job_name = "tf_test_%s"
	tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
  }

resource "tencentcloud_dts_migrate_job" "job" {
  	service_id = tencentcloud_dts_migrate_service.service.id
	run_mode = "immediate"
	migrate_option {
		database_table {
			object_mode = "partial"
			databases {
				db_name = "tf_ci_test"
				db_mode = "partial"
				table_mode = "partial"
				// view_mode = "partial"
				tables {
					table_name = "test"
					new_table_name = "test_%s"
					table_edit_mode = "rename"
				}
			}
		}
	}
	src_info {
			region = "ap-guangzhou"
			access_type = "cdb"
			database_type = "mysql"
			node_type = "simple"
			info {
				user = "keep_dts"
				password = "Letmein123"
				instance_id = "cdb-fitq5t9h"
			}

	}
	dst_info {
			region = "ap-guangzhou"
			access_type = "cdb"
			database_type = "cynosdbmysql"
			node_type = "simple"
			info {
				user = "root"
				password = "cynos@123"
				instance_id = tencentcloud_cynosdb_cluster.foo.id
			}
	}
	auto_retry_time_range_minutes = 0
}

resource "tencentcloud_dts_migrate_job_start_operation" "start"{
	job_id = tencentcloud_dts_migrate_job.job.id
}

`
