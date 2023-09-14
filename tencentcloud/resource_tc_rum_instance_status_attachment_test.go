package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumInstanceStatusAttachmentResource_basic -v
func TestAccTencentCloudRumInstanceStatusAttachmentResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumInstanceStatusAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_instance_status_attachment.instance_status_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_rum_instance_status_attachment.instance_status_attachment", "instance_status", "6"),
				),
			},
			{
				Config: testAccRumInstanceStatusAttachmentUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_instance_status_attachment.instance_status_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_rum_instance_status_attachment.instance_status_attachment", "instance_status", "2"),
				),
			},
			{
				ResourceName:            "tencentcloud_rum_instance_status_attachment.instance_status_attachment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"operate"},
			},
		},
	})
}

const testAccRumInstanceStatusAttachmentVar = `
variable "instance_id" {
  default = "` + defaultRumInstanceId + `"
}
`

const testAccRumInstanceStatusAttachment = testAccRumInstanceStatusAttachmentVar + `

resource "tencentcloud_rum_instance_status_attachment" "instance_status_attachment" {
	instance_id = var.instance_id
	operate = "stop"
}

`

const testAccRumInstanceStatusAttachmentUp = testAccRumInstanceStatusAttachmentVar + `

resource "tencentcloud_rum_instance_status_attachment" "instance_status_attachment" {
	instance_id = var.instance_id
	operate = "resume"
}

`
