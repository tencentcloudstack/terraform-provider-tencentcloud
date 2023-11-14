package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbFileDownloadUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbFileDownloadUrlDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_file_download_url.file_download_url")),
			},
		},
	})
}

const testAccDcdbFileDownloadUrlDataSource = `

data "tencentcloud_dcdb_file_download_url" "file_download_url" {
  instance_id = ""
  shard_id = ""
  file_path = ""
  }

`
