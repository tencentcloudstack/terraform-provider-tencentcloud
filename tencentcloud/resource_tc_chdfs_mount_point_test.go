package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudChdfsMountPointResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsMountPoint,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_chdfs_mount_point.mount_point", "id")),
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
  mount_point_name = &lt;nil&gt;
  file_system_id = &lt;nil&gt;
  mount_point_status = &lt;nil&gt;
}

`
