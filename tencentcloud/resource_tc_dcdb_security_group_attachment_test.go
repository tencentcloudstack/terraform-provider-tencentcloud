package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbSecurityGroupAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbSecurityGroupAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_security_group_attachment.security_group_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_security_group_attachment.security_group_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbSecurityGroupAttachment = `

resource "tencentcloud_dcdb_security_group_attachment" "security_group_attachment" {
  security_group_id = &lt;nil&gt;
  instance_id = &lt;nil&gt;
}

`
