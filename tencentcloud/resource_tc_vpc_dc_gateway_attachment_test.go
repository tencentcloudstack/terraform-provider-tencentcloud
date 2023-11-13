package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcDcGatewayAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcDcGatewayAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_dc_gateway_attachment.dc_gateway_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_dc_gateway_attachment.dc_gateway_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcDcGatewayAttachment = `

resource "tencentcloud_vpc_dc_gateway_attachment" "dc_gateway_attachment" {
  vpc_id = "vpc-111"
  nat_gateway_id = "nat-test123"
  direct_connect_gateway_id = "dcg-test123"
}

`
