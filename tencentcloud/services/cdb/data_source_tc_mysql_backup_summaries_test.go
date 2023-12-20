package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlBackupSummariesDataSource_basic -v
func TestAccTencentCloudMysqlBackupSummariesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlBackupSummariesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_backup_summaries.backup_summaries"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_summaries.backup_summaries", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_summaries.backup_summaries", "items.0.auto_backup_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_summaries.backup_summaries", "items.0.auto_backup_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_summaries.backup_summaries", "items.0.backup_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_summaries.backup_summaries", "items.0.binlog_backup_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_summaries.backup_summaries", "items.0.binlog_backup_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_summaries.backup_summaries", "items.0.data_backup_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_summaries.backup_summaries", "items.0.data_backup_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_summaries.backup_summaries", "items.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_summaries.backup_summaries", "items.0.manual_backup_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_summaries.backup_summaries", "items.0.manual_backup_volume"),
				),
			},
		},
	})
}

const testAccMysqlBackupSummariesDataSource = `

data "tencentcloud_mysql_backup_summaries" "backup_summaries" {
  product = "mysql"
  order_by = "BackupVolume"
  order_direction = "ASC"
}

`
