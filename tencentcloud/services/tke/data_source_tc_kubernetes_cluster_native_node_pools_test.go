package tke

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
			),
		}},
	})
}

const testAccKubernetesClusterNativeNodePoolsDataSource = `

data "tencentcloud_kubernetes_cluster_native_node_pools" "kubernetes_cluster_native_node_pools" {
  filters = {
  }
}
`
