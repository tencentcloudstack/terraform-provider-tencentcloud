package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcGatewayFlowQosDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcGatewayFlowQosDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_gateway_flow_qos.gateway_flow_qos")),
			},
		},
	})
}

const testAccVpcGatewayFlowQosDataSource = `

data "tencentcloud_vpc_gateway_flow_qos" "gateway_flow_qos" {
  gateway_id = "vpngw-gt8bianl"
}

`
