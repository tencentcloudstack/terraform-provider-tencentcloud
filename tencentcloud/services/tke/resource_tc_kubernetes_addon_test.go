package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesAddonResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAddon,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon.kubernetes_addon", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon.kubernetes_addon", "addon_name", "cos"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon.kubernetes_addon", "addon_version", "2018-05-25"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon.kubernetes_addon", "phase"),
					// resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon.kubernetes_addon", "reason"),
				),
			},
			{
				Config: testAccKubernetesAddonUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon.kubernetes_addon", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon.kubernetes_addon", "addon_name", "cos"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_addon.kubernetes_addon", "addon_version", "2018-05-25"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon.kubernetes_addon", "phase"),
					// resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_addon.kubernetes_addon", "reason"),
				),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_addon.kubernetes_addon",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKubernetesAddon = `
// resource "tencentcloud_kubernetes_cluster" "example" {
//   vpc_id                  = "` + tcacctest.DefaultTmpVpcId + `"
//   cluster_cidr            = "10.31.0.0/16"
//   cluster_max_pod_num     = 32
//   cluster_name            = "tf_example_cluster"
//   cluster_desc            = "example for tke cluster"
//   cluster_max_service_num = 32
//   cluster_internet        = false # (can be ignored) open it after the nodes added
//   cluster_version         = "1.22.5"
//   cluster_deploy_type     = "MANAGED_CLUSTER"
//   # without any worker config
// }

resource "tencentcloud_kubernetes_addon" "kubernetes_addon" {
  # cluster_id = tencentcloud_kubernetes_cluster.example.id
  cluster_id = "cls-lv0y4v68"
  addon_name    = "cos"
  addon_version = "2018-05-25"
  raw_values    = "{}"
}
`

const testAccKubernetesAddonUpdate = `
resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                  = "` + tcacctest.DefaultTmpVpcId + `"
  cluster_cidr            = "10.31.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "tf_example_cluster" 
  cluster_desc            = "example for tke cluster"
  cluster_max_service_num = 32
  cluster_internet        = false # (can be ignored) open it after the nodes added
  cluster_version         = "1.22.5"
  cluster_deploy_type     = "MANAGED_CLUSTER"
  # without any worker config
}

resource "tencentcloud_kubernetes_addon" "kubernetes_addon" {
  cluster_id  = tencentcloud_kubernetes_cluster.example.id
  addon_name    = "cos"
  addon_version = "2018-05-25"
  raw_values    = "{\"tolerations\":[{\"key\":\"test\",\"value\":\"100\",\"operator\":\"Equal\"}]}"
}
`
