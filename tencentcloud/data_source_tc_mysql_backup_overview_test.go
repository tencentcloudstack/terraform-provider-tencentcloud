package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlBackupOverviewDataSource_basic -v
func TestAccTencentCloudMysqlBackupOverviewDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlBackupOverviewDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_backup_overview.backup_overview"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_overview.backup_overview", "backup_archive_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_overview.backup_overview", "backup_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_overview.backup_overview", "backup_standby_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_overview.backup_overview", "backup_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_overview.backup_overview", "billing_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_overview.backup_overview", "free_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_overview.backup_overview", "product"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_backup_overview.backup_overview", "remote_backup_volume"),
				),
			},
		},
	})
}

const testAccMysqlBackupOverviewDataSource = `

data "tencentcloud_mysql_backup_overview" "backup_overview" {
  product = "mysql"
}

`
