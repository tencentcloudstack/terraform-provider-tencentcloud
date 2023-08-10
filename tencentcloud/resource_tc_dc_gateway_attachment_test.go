package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcGatewayAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcGatewayAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dc_gateway_attachment.dc_gateway_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_dc_gateway_attachment.dc_gateway_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcGatewayAttachment = `
resource "tencentcloud_vpc" "example" {
  name = "tf-vpc"
  cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_dc_gateway_attachment" "dc_gateway_attachment" {
  vpc_id = tencentcloud_vpc.example.id
  nat_gateway_id = "nat-7kanjc6y"
  direct_connect_gateway_id = "dcg-dmbhf7jf"
}

`
