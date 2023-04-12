package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlDataBackupOverviewDataSource_basic -v
func TestAccTencentCloudMysqlDataBackupOverviewDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDataBackupOverviewDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_data_backup_overview.data_backup_overview"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "auto_backup_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "auto_backup_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "data_backup_archive_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "data_backup_archive_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "data_backup_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "data_backup_standby_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "data_backup_standby_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "data_backup_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "manual_backup_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "manual_backup_volume"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "product"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "remote_backup_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_data_backup_overview.data_backup_overview", "remote_backup_volume"),
				),
			},
		},
	})
}

const testAccMysqlDataBackupOverviewDataSource = `

data "tencentcloud_mysql_data_backup_overview" "data_backup_overview" {
  product = "mysql"
                        }

`
