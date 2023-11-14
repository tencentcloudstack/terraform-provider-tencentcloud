package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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

resource "tencentcloud_chdfs_mount_point_attachment" "mount_point_attachment" {
  mount_point_id = &lt;nil&gt;
  access_group_ids = &lt;nil&gt;
}

`
