package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbBackupSummariesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbBackupSummariesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_backup_summaries.backup_summaries")),
			},
		},
	})
}

const testAccCdbBackupSummariesDataSource = `

data "tencentcloud_cdb_backup_summaries" "backup_summaries" {
  product = "mysql"
  order_by = "BackupVolume"
  order_direction = "ASC"
  }

`
