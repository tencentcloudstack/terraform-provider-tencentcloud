package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcGatewayFlowQosDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcGatewayFlowQosDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_gateway_flow_qos.gateway_flow_qos")),
			},
		},
	})
}

const testAccVpcGatewayFlowQosDataSource = `

data "tencentcloud_vpc_gateway_flow_qos" "gateway_flow_qos" {
  gateway_id = "vpngw-gt8bianl"
}

`
