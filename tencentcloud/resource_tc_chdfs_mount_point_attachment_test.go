package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudChdfsMountPointAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsMountPointAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_chdfs_mount_point_attachment.mount_point_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_chdfs_mount_point_attachment.mount_point_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccChdfsMountPointAttachment = `

resource "tencentcloud_chdfs_mount_point" "mount_point" {
  file_system_id     = "f14mpfy5lh4e"
  mount_point_name   = "terraform-test-mount-attach"
  mount_point_status = 1
}

resource "tencentcloud_chdfs_mount_point_attachment" "mount_point_attachment" {
  access_group_ids = [
    "ag-bvmzrbsm",
  ]
  mount_point_id   = tencentcloud_chdfs_mount_point.mount_point.id
}

`
