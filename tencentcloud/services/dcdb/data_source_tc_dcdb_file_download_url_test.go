package dcdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbNeedFixFileDownloadUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbFileDownloadUrlDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_file_download_url.file_download_url"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_log_files.log_files", "shard_id", "shard-1b5r04az"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_log_files.log_files", "file_path"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_log_files.log_files", "pre_signed_url"),
				),
			},
		},
	})
}

const testAccDcdbFileDownloadUrlDataSource = tcacctest.CommonPresetDcdb + `

data "tencentcloud_dcdb_file_download_url" "file_download_url" {
  instance_id = local.dcdb_id
  shard_id = "shard-1b5r04az"
  file_path = "/cos_backup/test.txt"
}
`
