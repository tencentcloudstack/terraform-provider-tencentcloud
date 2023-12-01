package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesClusterInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_cluster_instances.cluster_instances"),
					resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_cluster_instances.cluster_instances", "instance_set.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_cluster_instances.cluster_instances", "instance_set.0.instance_id", "ins-kqmx8dm2"),
					resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_cluster_instances.cluster_instances", "instance_set.0.instance_role", "WORKER"),
					resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_cluster_instances.cluster_instances", "instance_set.0.instance_state", "running"),
				),
			},
			{
				Config: testAccKubernetesClusterInstancesDataSourceFilter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_cluster_instances.cluster_instances"),
					resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_cluster_instances.cluster_instances", "instance_set.#", "0"),
				),
			},
		},
	})
}

const testAccKubernetesClusterInstancesDataSource = `
data "tencentcloud_kubernetes_cluster_instances" "cluster_instances" {
  cluster_id    = "cls-ely08ic4"
  instance_ids  = ["ins-kqmx8dm2"]
  instance_role = "WORKER"
}
`

const testAccKubernetesClusterInstancesDataSourceFilter = `
data "tencentcloud_kubernetes_cluster_instances" "cluster_instances" {
  cluster_id    = "cls-ely08ic4"
  instance_ids  = ["ins-kqmx8dm2"]
  instance_role = "WORKER"
  filters {
    name   = "nodepool-id"
    values = ["np-p4e6whqu"]
  }
}
`
