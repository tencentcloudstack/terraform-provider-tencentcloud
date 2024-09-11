package tke

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesNativeNodePoolsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccKubernetesNativeNodePools,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_native_node_pools.kubernetes_native_node_pools", "id")),
		}, {
			ResourceName:      "tencentcloud_kubernetes_native_node_pools.kubernetes_native_node_pools",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccKubernetesNativeNodePools = `

resource "tencentcloud_kubernetes_native_node_pools" "kubernetes_native_node_pools" {
  labels = {
  }
  taints = {
  }
  tags = {
    tags = {
    }
  }
  native = {
    system_disk = {
    }
    data_disks = {
    }
  }
}
`
