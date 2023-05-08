package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlBinlogBackupOverviewDataSource_basic -v
func TestAccTencentCloudMysqlBinlogBackupOverviewDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlBinlogBackupOverviewDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_binlog_backup_overview.binlog_backup_overview"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_binlog_backup_overview.binlog_backup_overview", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_binlog_backup_overview.binlog_backup_overview", "binlog_archive_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_binlog_backup_overview.binlog_backup_overview", "binlog_archive_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_binlog_backup_overview.binlog_backup_overview", "binlog_backup_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_binlog_backup_overview.binlog_backup_overview", "binlog_backup_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_binlog_backup_overview.binlog_backup_overview", "binlog_standby_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_binlog_backup_overview.binlog_backup_overview", "binlog_standby_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_binlog_backup_overview.binlog_backup_overview", "product"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_binlog_backup_overview.binlog_backup_overview", "remote_binlog_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_binlog_backup_overview.binlog_backup_overview", "remote_binlog_volume"),
				),
			},
		},
	})
}

const testAccMysqlBinlogBackupOverviewDataSource = `

data "tencentcloud_mysql_binlog_backup_overview" "binlog_backup_overview" {
  product = "mysql"
}

`
