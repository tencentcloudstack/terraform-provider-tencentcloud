package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesClusterRollOutSequenceTagConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterRollOutSequenceTagConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example", "tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example", "tags.0.key", "Env"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example", "tags.0.value", "Test"),
				),
			},
			{
				Config: testAccKubernetesClusterRollOutSequenceTagConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example", "tags.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example", "tags.0.key", "Env"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example", "tags.0.value", "Production"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example", "tags.1.key", "Protection-Level"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example", "tags.1.value", "High"),
				),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKubernetesClusterRollOutSequenceTagConfig = `
resource "tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config" "example" {
  cluster_id = "cls-xxxxxxxx"

  tags {
    key   = "Env"
    value = "Test"
  }
}
`

const testAccKubernetesClusterRollOutSequenceTagConfigUpdate = `
resource "tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config" "example" {
  cluster_id = "cls-xxxxxxxx"

  tags {
    key   = "Env"
    value = "Production"
  }

  tags {
    key   = "Protection-Level"
    value = "High"
  }
}
`
