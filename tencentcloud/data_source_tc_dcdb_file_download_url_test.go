package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbNeedFixFileDownloadUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbFileDownloadUrlDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_file_download_url.file_download_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_file_download_url.file_download_url", "shard_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_file_download_url.file_download_url", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_file_download_url.file_download_url", "file_path"),
				),
			},
		},
	})
}

const testAccDcdbFileDownloadUrlDataSource = CommonPresetDcdb + `

data "tencentcloud_dcdb_file_download_url" "file_download_url" {
  instance_id = local.dcdb_id
  shard_id = "shard-l86azfrj"
  file_path = "cos_backup/tdsql/group_1667305846_5423296/set_1686219511_15/xtrabackup/2023-08-04/cos_xtrabackup+1691111110+20230804+090510+3480354987+xbstream.lz4"
}
`
