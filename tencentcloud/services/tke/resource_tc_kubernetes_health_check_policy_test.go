package tke

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesHealthCheckPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccKubernetesHealthCheckPolicy,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "id"),
				resource.TestCheckResourceAttr("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "health_check_policy"),
			),
		}, {
			ResourceName:      "tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccKubernetesHealthCheckPolicy = `

resource "tencentcloud_kubernetes_health_check_policy" "kubernetes_health_check_policy" {
  cluster_id = "cls-eh0da110"
  health_check_policy = {
    name = "NP1"
    rules = {
		name = "RuntimeUnhealthy"
		auto_repair_enabled = true
		enabled	= true
    }
  }
}
`
