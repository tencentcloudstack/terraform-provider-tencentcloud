package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcEniSgAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEniSgAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_eni_sg_attachment.eni_sg_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_eni_sg_attachment.eni_sg_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcEniSgAttachment = `

resource "tencentcloud_eni_sg_attachment" "eni_sg_attachment" {
  network_interface_ids = ["eni-p0hkgx8p"]
  security_group_ids    = ["sg-902tl7t7", "sg-edmur627"]
}

`
