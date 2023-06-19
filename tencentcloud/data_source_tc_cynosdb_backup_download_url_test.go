package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCynosdbBackupDownloadUrlDataSource_basic -v
func TestAccTencentCloudNeedFixCynosdbBackupDownloadUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbBackupDownloadUrlDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_backup_download_url.backup_download_url"),
				),
			},
		},
	})
}

const testAccCynosdbBackupDownloadUrlDataSource = `
data "tencentcloud_cynosdb_backup_download_url" "backup_download_url" {
  cluster_id = "cynosdbmysql-bws8h88b"
  backup_id  = 480782
}
`
