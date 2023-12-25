package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixVpcLocalGatewayResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcLocalGateway,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_local_gateway.local_gateway", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_local_gateway.local_gateway",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcLocalGateway = `

resource "tencentcloud_vpc_local_gateway" "local_gateway" {
  local_gateway_name = "local-gw-test"
  vpc_id             = "vpc-lh4nqig9"
  cdc_id             = "cluster-j9gyu1iy"
}

`
