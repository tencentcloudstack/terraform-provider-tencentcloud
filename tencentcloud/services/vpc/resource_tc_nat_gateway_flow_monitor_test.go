package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNatGatewayFlowMonitorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayFlowMonitor,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_nat_gateway_flow_monitor.example", "id")),
			},
		},
	})
}

const testAccNatGatewayFlowMonitor = `

resource "tencentcloud_nat_gateway_flow_monitor" "example" {
  gateway_id = "nat-e6u6axsm"
  enable     = true
}

`
