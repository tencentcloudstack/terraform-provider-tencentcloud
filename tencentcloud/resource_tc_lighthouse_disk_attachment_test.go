package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudLighthouseDiskAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
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

const testAccLighthouseDiskAttachment = DefaultLighthoustVariables + `
resource "tencentcloud_lighthouse_disk_attachment" "disk_attachment" {
	disk_id = var.lighthouse_disk_id
	instance_id = var.lighthouse_id
}
`
