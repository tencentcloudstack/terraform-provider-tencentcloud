package clb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTestingClbSecurityGroupAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
