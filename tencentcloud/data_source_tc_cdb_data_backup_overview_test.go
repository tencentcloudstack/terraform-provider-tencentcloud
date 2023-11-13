package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbDataBackupOverviewDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbDataBackupOverviewDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_data_backup_overview.data_backup_overview")),
			},
		},
	})
}

const testAccCdbDataBackupOverviewDataSource = `

data "tencentcloud_cdb_data_backup_overview" "data_backup_overview" {
  product = "mysql"
                        }

`
