package tke_test

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
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesHealthCheckPolicyCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "name", "example"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "rules.0.auto_repair_enabled", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "rules.0.enabled", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "rules.1.auto_repair_enabled", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "rules.1.enabled", "true"),
				),
			}, {
				ResourceName:      "tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy",
				ImportState:       true,
				ImportStateVerify: true,
			}, {
				Config: testAccKubernetesHealthCheckPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "name", "example"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "rules.0.auto_repair_enabled", "false"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "rules.0.enabled", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "rules.1.auto_repair_enabled", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy", "rules.1.enabled", "false"),
				),
			}},
	})
}

// const testAccKubernetesHealthCheckPolicyCreate = testAccTkeCluster + `
const testAccKubernetesHealthCheckPolicyCreate = `

resource "tencentcloud_kubernetes_health_check_policy" "kubernetes_health_check_policy" {
	cluster_id = "cls-eh0da110"

	name = "example"
	rules {
		name = "OOMKilling"
		auto_repair_enabled = true
		enabled	= true
	}
	rules {
		name = "KubeletUnhealthy"
		auto_repair_enabled = true
		enabled	= true
	}
}
`

// const testAccKubernetesHealthCheckPolicyUpdate = testAccTkeCluster + `
const testAccKubernetesHealthCheckPolicyUpdate = `

resource "tencentcloud_kubernetes_health_check_policy" "kubernetes_health_check_policy" {
	cluster_id = "cls-eh0da110"
	name = "example"
	rules {
		name = "OOMKilling"
		auto_repair_enabled = false
		enabled	= true
	}
	rules {
		name = "KubeletUnhealthy"
		auto_repair_enabled = true
		enabled	= false
	}
}
`
