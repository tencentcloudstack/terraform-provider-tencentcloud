package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesAddonConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAddonConfig,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon_config.kubernetes_addon_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon_config.kubernetes_addon_config", "addon_name", "cluster-autoscaler"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon_config.kubernetes_addon_config", "raw_values", "{\"extraArgs\":{\"scale-down-enabled\":true,\"max-empty-bulk-delete\":11,\"scale-down-delay-after-add\":\"10mm\",\"scale-down-unneeded-time\":\"10mm\",\"scale-down-utilization-threshold\":0.005,\"ignore-daemonsets-utilization\":false,\"skip-nodes-with-local-storage\":true,\"skip-nodes-with-system-pods\":true}}"),
				),
			},
			{
				Config: testAccKubernetesAddonConfigUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon_config.kubernetes_addon_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon_config.kubernetes_addon_config", "addon_name", "cluster-autoscaler"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon_config.kubernetes_addon_config", "raw_values", "{\"extraArgs\":{\"scale-down-enabled\":false,\"max-empty-bulk-delete\":10,\"scale-down-delay-after-add\":\"10mm\",\"scale-down-unneeded-time\":\"10mm\",\"scale-down-utilization-threshold\":0.005,\"ignore-daemonsets-utilization\":false,\"skip-nodes-with-local-storage\":true,\"skip-nodes-with-system-pods\":true}}"),
				),
			},
		},
	})
}

const testAccKubernetesAddonConfig = `
resource "tencentcloud_kubernetes_addon_config" "kubernetes_addon_config" {
	cluster_id = "cls-bzoq8t02"
	addon_name = "cluster-autoscaler"
	raw_values = "{\"extraArgs\":{\"scale-down-enabled\":true,\"max-empty-bulk-delete\":11,\"scale-down-delay-after-add\":\"10mm\",\"scale-down-unneeded-time\":\"10mm\",\"scale-down-utilization-threshold\":0.005,\"ignore-daemonsets-utilization\":false,\"skip-nodes-with-local-storage\":true,\"skip-nodes-with-system-pods\":true}}"
}
`
const testAccKubernetesAddonConfigUpdate = `
resource "tencentcloud_kubernetes_addon_config" "kubernetes_addon_config" {
	cluster_id = "cls-bzoq8t02"
	addon_name = "cluster-autoscaler"
	raw_values = "{\"extraArgs\":{\"scale-down-enabled\":false,\"max-empty-bulk-delete\":10,\"scale-down-delay-after-add\":\"10mm\",\"scale-down-unneeded-time\":\"10mm\",\"scale-down-utilization-threshold\":0.005,\"ignore-daemonsets-utilization\":false,\"skip-nodes-with-local-storage\":true,\"skip-nodes-with-system-pods\":true}}"
}
`
