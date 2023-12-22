package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixMariadbFileDownloadUrlDataSource_basic -v
func TestAccTencentCloudNeedFixMariadbFileDownloadUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbFileDownloadUrlDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_file_download_url.file_download_url"),
				),
			},
		},
	})
}

const testAccMariadbFileDownloadUrlDataSource = `
data "tencentcloud_mariadb_file_download_url" "file_download_url" {
  instance_id = "tdsql-9vqvls95"
  file_path   = "/cos_backup/test.txt"
}
`
