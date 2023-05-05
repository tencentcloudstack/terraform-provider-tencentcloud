package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudChdfsMountPointsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsMountPointsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_chdfs_mount_points.mount_points")),
			},
		},
	})
}

const testAccChdfsMountPointsDataSource = `

data "tencentcloud_chdfs_mount_points" "mount_points" {
  file_system_id     = "f14mpfy5lh4e"
}

`
