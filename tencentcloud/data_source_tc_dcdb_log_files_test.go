package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_log_files.log_files")),
			},
		},
	})
}

const testAccDcdbLogFilesDataSource = `

data "tencentcloud_dcdb_log_files" "log_files" {
  instance_id = &lt;nil&gt;
  shard_id = &lt;nil&gt;
  type = &lt;nil&gt;
      }

`
