package chdfs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudChdfsMountPointResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsMountPoint,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_chdfs_mount_point.mount_point", "id")),
			},
			{
				Config: testAccChdfsMountPointUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_chdfs_mount_point.mount_point", "id"),
					resource.TestCheckResourceAttr("tencentcloud_chdfs_mount_point.mount_point", "mount_point_name", "terraform-for-test"),
					resource.TestCheckResourceAttr("tencentcloud_chdfs_mount_point.mount_point", "mount_point_status", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_chdfs_mount_point.mount_point",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccChdfsMountPoint = `

resource "tencentcloud_chdfs_mount_point" "mount_point" {
  file_system_id     = "f14mpfy5lh4e"
  mount_point_name   = "terraform-test"
  mount_point_status = 1
}

`

const testAccChdfsMountPointUpdate = `

resource "tencentcloud_chdfs_mount_point" "mount_point" {
  file_system_id     = "f14mpfy5lh4e"
  mount_point_name   = "terraform-for-test"
  mount_point_status = 2
}

`
