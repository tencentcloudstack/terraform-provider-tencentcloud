package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbFileDownloadUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbFileDownloadUrlDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_file_download_url.file_download_url")),
			},
		},
	})
}

const testAccMariadbFileDownloadUrlDataSource = `

data "tencentcloud_mariadb_file_download_url" "file_download_url" {
  instance_id = ""
  file_path = ""
  }

`
