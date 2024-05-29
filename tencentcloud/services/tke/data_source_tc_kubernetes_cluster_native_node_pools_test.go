package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesClusterNativeNodePoolsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccKubernetesClusterNativeNodePoolsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_cluster_native_node_pools.kubernetes_cluster_native_node_pools"),
				resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_cluster_native_node_pools.kubernetes_cluster_native_node_pools", "node_pools.#", "1"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_native_node_pools.kubernetes_cluster_native_node_pools", "node_pools.0.node_pool_id"),
				resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_cluster_native_node_pools.kubernetes_cluster_native_node_pools", "node_pools.0.name", "tf-native-node-pool1"),
			),
		}},
	})
}

const testAccKubernetesClusterNativeNodePoolsDataSource = testAccTkeNativeNodePool + `
data "tencentcloud_kubernetes_cluster_native_node_pools" "kubernetes_cluster_native_node_pools" {
  cluster_id = tencentcloud_kubernetes_cluster.kubernetes_cluster.id
  filters {
    name   = "NodePoolsName"
    values = [tencentcloud_kubernetes_native_node_pool.native_node_pool_test.name]
  }
}
`
