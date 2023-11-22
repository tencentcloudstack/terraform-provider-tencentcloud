package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTestingClbSecurityGroupAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingClbSecurityGroupAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_security_group_attachment.security_group_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_clb_security_group_attachment.security_group_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTestingClbSecurityGroupAttachment = `

resource "tencentcloud_clb_security_group_attachment" "security_group_attachment" {
  security_group = "sg-46x20487"
  load_balancer_ids = ["lb-d72hklcc"]
}

`
