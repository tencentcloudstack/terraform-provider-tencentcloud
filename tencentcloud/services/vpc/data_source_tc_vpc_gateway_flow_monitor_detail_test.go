package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcGatewayFlowMonitorDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcGatewayFlowMonitorDetailDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_gateway_flow_monitor_detail.gateway_flow_monitor_detail")),
			},
		},
	})
}

const testAccVpcGatewayFlowMonitorDetailDataSource = `

data "tencentcloud_vpc_gateway_flow_monitor_detail" "gateway_flow_monitor_detail" {
  time_point      = "2023-06-02 12:15:20"
  vpn_id          = "vpngw-gt8bianl"
  order_field     = "OutTraffic"
  order_direction = "DESC"
}

`
