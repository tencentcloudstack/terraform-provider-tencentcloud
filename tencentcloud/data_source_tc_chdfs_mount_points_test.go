package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
  file_system_id = &lt;nil&gt;
  access_group_id = &lt;nil&gt;
  owner_uin = &lt;nil&gt;
  }

`
