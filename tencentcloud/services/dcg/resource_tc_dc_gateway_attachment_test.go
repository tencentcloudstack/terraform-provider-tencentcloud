package dcg_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcGatewayAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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

resource "tencentcloud_dc_gateway_attachment" "dc_gateway_attachment" {
  vpc_id = "vpc-4h9v4mo3"
  nat_gateway_id = "nat-7kanjc6y"
  direct_connect_gateway_id = "dcg-dmbhf7jf"
}

`
