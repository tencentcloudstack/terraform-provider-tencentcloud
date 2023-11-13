package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseGatewayNodesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGatewayNodesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateway_nodes.gateway_nodes")),
			},
		},
	})
}

const testAccTseGatewayNodesDataSource = `

data "tencentcloud_tse_gateway_nodes" "gateway_nodes" {
  gateway_id = "gateway-xx"
  group_id = "group-xx"
  }

`
