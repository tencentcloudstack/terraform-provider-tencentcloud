package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixVpcClassicLinkAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcClassicLinkAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_classic_link_attachment.classic_link_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_classic_link_attachment.classic_link_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcClassicLinkAttachment = `

resource "tencentcloud_vpc_classic_link_attachment" "classic_link_attachment" {
  vpc_id       = "vpc-hdvfe0g1"
  instance_ids = ["ins-ceynqvnu"]
}

`
