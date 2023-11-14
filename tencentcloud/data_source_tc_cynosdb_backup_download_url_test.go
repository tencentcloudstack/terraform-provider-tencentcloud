package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbBackupDownloadUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbBackupDownloadUrlDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_backup_download_url.backup_download_url")),
			},
		},
	})
}

const testAccCynosdbBackupDownloadUrlDataSource = `

data "tencentcloud_cynosdb_backup_download_url" "backup_download_url" {
  cluster_id = "cynosdbmysql-123"
  backup_id = 100
  }

`
