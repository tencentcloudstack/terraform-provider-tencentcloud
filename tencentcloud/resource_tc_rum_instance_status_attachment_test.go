package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRumInstanceStatusAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumInstanceStatusAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_rum_instance_status_attachment.instance_status_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_rum_instance_status_attachment.instance_status_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRumInstanceStatusAttachment = `

resource "tencentcloud_rum_instance_status_attachment" "instance_status_attachment" {
  instance_id = "rum-xxx"
}

`
