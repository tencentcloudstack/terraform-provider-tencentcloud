package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbLogFilesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbLogFilesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_log_files.log_files"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_log_files.log_files", "shard_id", "shard-1b5r04az"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_log_files.log_files", "type", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_log_files.log_files", "files.#"),
				),
			},
		},
	})
}

const testAccDcdbLogFilesDataSource = CommonPresetDcdb + `

data "tencentcloud_dcdb_log_files" "log_files" {
	instance_id = local.dcdb_id
	shard_id    = "shard-1b5r04az"
	type        = 1
}

`
