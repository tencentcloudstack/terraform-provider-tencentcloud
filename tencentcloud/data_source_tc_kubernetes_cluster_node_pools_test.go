package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTkeClusterNodePoolsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeClusterNodePoolsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_cluster_node_pools.cluster_node_pools"),
					resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_cluster_node_pools.cluster_node_pools", "node_pool_set.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_cluster_node_pools.cluster_node_pools", "node_pool_set.0.node_pool_id", "np-ngjwhdv4"),
					resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_cluster_node_pools.cluster_node_pools", "node_pool_set.0.name", "mynodepool_xxxx"),
				),
			},
		},
	})
}

const testAccTkeClusterNodePoolsDataSource = `
data "tencentcloud_kubernetes_cluster_node_pools" "cluster_node_pools" {
  cluster_id = "cls-kzilgv5m"
  filters {
    name   = "NodePoolsName"
    values = ["mynodepool_xxxx"]
  }
  filters {
    name   = "NodePoolsId"
    values = ["np-ngjwhdv4"]
  }
}
`
