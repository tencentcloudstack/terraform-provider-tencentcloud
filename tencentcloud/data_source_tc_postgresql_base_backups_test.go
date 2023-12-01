package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPostgresqlBaseBackupsObject = "data.tencentcloud_postgresql_base_backups.base_backups"

func TestAccTencentCloudPostgresqlBaseBackupsDataSource_basic(t *testing.T) {
	// t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -7).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().AddDate(0, 0, 1).In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccStepSetRegion(t, "ap-chengdu")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccPostgresqlBaseBackupsDataSource_bytime, startTime, endTime),
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-chengdu")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testAccPostgresqlBaseBackupsObject),
					resource.TestCheckResourceAttr(testAccPostgresqlBaseBackupsObject, "min_finish_time", startTime),
					resource.TestCheckResourceAttr(testAccPostgresqlBaseBackupsObject, "max_finish_time", endTime),
					resource.TestCheckResourceAttr(testAccPostgresqlBaseBackupsObject, "order_by", "StartTime"),
					resource.TestCheckResourceAttr(testAccPostgresqlBaseBackupsObject, "order_by_type", "asc"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.#"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.name"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.backup_method"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.backup_mode"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.state"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.size"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.start_time"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.finish_time"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.expire_time"),
				),
			},
			{
				Config: testAccPostgresqlBaseBackupsDataSource_byMasterId,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-chengdu")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testAccPostgresqlBaseBackupsObject),
					resource.TestCheckResourceAttr(testAccPostgresqlBaseBackupsObject, "filters.0.name", "db-instance-id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "filters.0.values.#"),
					resource.TestCheckResourceAttr(testAccPostgresqlBaseBackupsObject, "order_by", "Size"),
					resource.TestCheckResourceAttr(testAccPostgresqlBaseBackupsObject, "order_by_type", "asc"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.#"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.name"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.backup_method"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.backup_mode"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.state"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.size"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.start_time"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.finish_time"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.expire_time"),
				),
			},
			{
				Config: testAccPostgresqlBaseBackupsDataSource_byBackupId,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-chengdu")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testAccPostgresqlBaseBackupsObject),
					resource.TestCheckResourceAttr(testAccPostgresqlBaseBackupsObject, "filters.0.name", "base-backup-id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "filters.0.values.#"),
					resource.TestCheckResourceAttr(testAccPostgresqlBaseBackupsObject, "order_by", "Size"),
					resource.TestCheckResourceAttr(testAccPostgresqlBaseBackupsObject, "order_by_type", "asc"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.#"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.name"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.backup_method"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.backup_mode"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.state"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.size"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.start_time"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.finish_time"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlBaseBackupsObject, "base_backup_set.0.expire_time"),
				),
			},
		},
	})
}

const testAccPostgresqlBaseBackupsDataSource_bytime = testAccPostgresqlInstance_common + `

data "tencentcloud_postgresql_base_backups" "base_backups" {
  min_finish_time = "%s"
  max_finish_time = "%s"

  order_by = "StartTime"
  order_by_type = "asc"
}

`

const testAccPostgresqlBaseBackupsDataSource_byMasterId = testAccPostgresqlInstance_common + `
data "tencentcloud_postgresql_base_backups" "base_backups" {
  filters {
    name   = "db-instance-id"
    values = [local.pgsql_id]
  }

  order_by      = "Size"
  order_by_type = "asc"
}


`

const testAccPostgresqlBaseBackupsDataSource_byBackupId = testAccPostgresqlInstance_common + `
data "tencentcloud_postgresql_base_backups" "query_master" {
  filters {
    name   = "db-instance-id"
    values = [local.pgsql_id]
  }

  order_by      = "Size"
  order_by_type = "asc"

}

data "tencentcloud_postgresql_base_backups" "base_backups" {
  filters {
    name   = "base-backup-id"
    values = [data.tencentcloud_postgresql_base_backups.query_master.base_backup_set.0.id]
  }

  order_by      = "Size"
  order_by_type = "asc"

}


`
