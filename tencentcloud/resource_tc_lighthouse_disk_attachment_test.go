package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseDiskAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseDiskAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_disk_attachment.disk_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_disk_attachment.disk_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseDiskAttachment = `

resource "tencentcloud_lighthouse_disk_attachment" "disk_attachment" {
  disk_ids = 
  instance_id = "lhins-123456"
  renew_flag = "NOTIFY_AND_MANUAL_RENEW"
}

`
