package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbBinlogBackupOverviewDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbBinlogBackupOverviewDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_binlog_backup_overview.binlog_backup_overview")),
			},
		},
	})
}

const testAccCdbBinlogBackupOverviewDataSource = `

data "tencentcloud_cdb_binlog_backup_overview" "binlog_backup_overview" {
  product = "mysql"
                }

`
