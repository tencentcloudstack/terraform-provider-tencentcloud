package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesClusterSchedulerPolicyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterSchedulerPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_scheduler_policy_config.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_scheduler_policy_config.example", "high_performance", "true"),
				),
			},
			{
				Config: testAccKubernetesClusterSchedulerPolicyConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_scheduler_policy_config.example", "high_performance", "false"),
				),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_cluster_scheduler_policy_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// NOTE: Replace cls-5e7wsn94 with a real cluster ID in your environment.
const testAccKubernetesClusterSchedulerPolicyConfig = `
resource "tencentcloud_kubernetes_cluster_scheduler_policy_config" "example" {
  cluster_id       = "cls-5e7wsn94"
  high_performance = true

  client_connection {
    qps   = 50
    burst = 100
  }
}
`

const testAccKubernetesClusterSchedulerPolicyConfigUpdate = `
resource "tencentcloud_kubernetes_cluster_scheduler_policy_config" "example" {
  cluster_id       = "cls-5e7wsn94"
  high_performance = false

  client_connection {
    qps   = 50
    burst = 100
  }
}
`
