package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesControlPlaneLogResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesControlPlaneLog,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_control_plane_log.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_control_plane_log.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKubernetesControlPlaneLog = `
resource "tencentcloud_kubernetes_control_plane_log" "example" {
  cluster_id   = "cls-rng1h5ei"
  cluster_type = "tke"
  components {
    name         = "cluster-autoscaler"
    log_level    = "2"
    log_set_id   = "40eed846-0f43-44b1-b216-c786a8970b1f"
    topic_id     = "21918a54-9ab4-40bc-90cd-c600cff00695"
    topic_region = "ap-guangzhou"
  }
}
`
