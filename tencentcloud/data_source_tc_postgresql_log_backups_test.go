package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlLogBackupsDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -7).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().AddDate(0, 0, 1).In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccPostgresqlLogBackupsDataSource, startTime, endTime),
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_log_backups.log_backups"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_log_backups.log_backups", "min_finish_time", startTime),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_log_backups.log_backups", "max_finish_time", endTime),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "filters.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_log_backups.log_backups", "filters.0.name", "db-instance-id"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_log_backups.log_backups", "filters.0.values.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_log_backups.log_backups", "order_by", "StartTime"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_log_backups.log_backups", "order_by_type", "desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "log_backup_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "log_backup_set.0.db_instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "log_backup_set.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "log_backup_set.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "log_backup_set.0.backup_method"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "log_backup_set.0.backup_mode"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "log_backup_set.0.state"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "log_backup_set.0.size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "log_backup_set.0.start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "log_backup_set.0.finish_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_log_backups.log_backups", "log_backup_set.0.expire_time"),
				),
			},
		},
	})
}

const testAccPostgresqlLogBackupsDataSource = CommonPresetPGSQL + `

data "tencentcloud_postgresql_log_backups" "log_backups" {
  min_finish_time = "%s"
  max_finish_time = "%s"
  filters {
		name = "db-instance-id"
		values = [local.pgsql_id]
  }
  order_by = "StartTime"
  order_by_type = "desc"
}
`
