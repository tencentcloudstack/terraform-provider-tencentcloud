package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTseGatewayNodesDataSource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_gateway_nodes.gateway_nodes"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_nodes.gateway_nodes", "node_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_nodes.gateway_nodes", "node_list.0.group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_nodes.gateway_nodes", "node_list.0.group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_nodes.gateway_nodes", "node_list.0.node_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_nodes.gateway_nodes", "node_list.0.node_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_nodes.gateway_nodes", "node_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_nodes.gateway_nodes", "node_list.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_gateway_nodes.gateway_nodes", "node_list.0.zone_id"),
				),
			},
		},
	})
}

const testAccTseGatewayNodesDataSource = `

data "tencentcloud_tse_gateway_nodes" "gateway_nodes" {
	gateway_id = "gateway-ddbb709b"
	group_id   = "group-013c0d8e"
}

`
