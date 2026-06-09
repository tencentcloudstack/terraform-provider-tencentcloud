package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesClusterExtraArgsConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterExtraArgsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_extra_args_config.example", "id"),
				),
			},
			{
				Config: testAccKubernetesClusterExtraArgsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_extra_args_config.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_cluster_extra_args_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// NOTE: Replace cls-xxxxxxxx with a real managed cluster ID in your environment.
const testAccKubernetesClusterExtraArgsConfig = `
resource "tencentcloud_kubernetes_cluster_extra_args_config" "example" {
  cluster_id = "cls-man1vvi2"
  kube_apiserver = [
    "goaway-chance=0",
    "kubelet-preferred-address-types=Hostname"
  ]
  kube_controller_manager = [
    "concurrent-serviceaccount-token-syncs=5"
  ]
  kube_scheduler = [
    "kube-api-qps=50"
  ]
}
`

const testAccKubernetesClusterExtraArgsConfigUpdate = `
resource "tencentcloud_kubernetes_cluster_extra_args_config" "example" {
  cluster_id = "cls-man1vvi2"
  kube_apiserver = [
    "goaway-chance=0",
    "kubelet-preferred-address-types=Hostname"
  ]
  kube_controller_manager = [
    "concurrent-serviceaccount-token-syncs=5"
  ]
  kube_scheduler = [
    "kube-api-qps=50"
  ]
}
`
