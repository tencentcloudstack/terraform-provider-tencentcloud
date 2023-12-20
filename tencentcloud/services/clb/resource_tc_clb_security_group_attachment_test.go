package clb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClbSecurityGroupAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbSecurityGroupAttachment,
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

const testAccClbSecurityGroupAttachment = `

resource "tencentcloud_clb_security_group_attachment" "security_group_attachment" {
  security_group = "sg-ijato2x1"
  load_balancer_ids = ["lb-5dnrkgry"]
}

`
