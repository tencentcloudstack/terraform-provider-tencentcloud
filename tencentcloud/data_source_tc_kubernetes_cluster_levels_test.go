package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudKubernetesClusterLevelDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDataSourceKubernetesClusterLevelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKubernetesClusterLevelBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_levels.foo", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_levels.foo", "list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_levels.foo", "list.0.alias"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_levels.foo", "list.0.crd_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_levels.foo", "list.0.config_map_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_levels.foo", "list.0.enable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_levels.foo", "list.0.node_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_levels.foo", "list.0.other_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_levels.foo", "list.0.pod_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_levels.with_cluster", "list.#"),
				),
			},
		},
	})
}

func testAccDataSourceKubernetesClusterLevelDestroy(s *terraform.State) error {
	return nil
}

const testAccDataSourceKubernetesClusterLevelBasic = `
data "tencentcloud_kubernetes_cluster_levels" "foo" {}

data "tencentcloud_kubernetes_clusters" "cls" {
  cluster_name = "` + defaultTkeClusterName + `"
}

data "tencentcloud_kubernetes_cluster_levels" "with_cluster" {
	cluster_id = data.tencentcloud_kubernetes_clusters.cls.list.0.cluster_id
}
`
