package cynosdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCynosdbBackupDownloadUrlDataSource_basic -v
func TestAccTencentCloudNeedFixCynosdbBackupDownloadUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbBackupDownloadUrlDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_backup_download_url.backup_download_url"),
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
