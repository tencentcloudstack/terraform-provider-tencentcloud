package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
				PreventDiskCleanup: true,
				Config:             fmt.Sprintf(testAccDtsMigrateJob_basic, curSec),
				Check:              resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job.job", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_migrate_job.job",
				ImportState:       true,
				ImportStateVerify: true,
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

const testAccDtsMigrateJob_basic = `

resource "tencentcloud_dts_migrate_service" "service" {
	src_database_type = "mysql"
	dst_database_type = "cynosdbmysql"
	src_region = "ap-guangzhou"
	dst_region = "ap-guangzhou"
	instance_class = "small"
	job_name = "tf_test_migration_service_1"
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
				user = "keep_dts"
				password = "Letmein123"
				instance_id = "cynosdbmysql-quqtcs13"
			}
	}
	job_name = "tf_migrate_job_config_test"
	auto_retry_time_range_minutes = 0
}

resource "tencentcloud_dts_migrate_job_start_operation" "start"{
	job_id = tencentcloud_dts_migrate_job.job.id
}

`
