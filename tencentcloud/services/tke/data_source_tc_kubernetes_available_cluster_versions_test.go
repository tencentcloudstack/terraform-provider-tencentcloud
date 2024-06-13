package tke_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudKubernetesAvailableClusterVersionsDataSource_basic -v
func TestAccTencentCloudKubernetesAvailableClusterVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAvailableClusterVersionsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_available_cluster_versions.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_available_cluster_versions.id", "versions.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_available_cluster_versions.id", "clusters.#"),
				),
			},
			{
				Config: testAccKubernetesAvailableClusterVersionsDataSource_multiple,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_available_cluster_versions.ids"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_available_cluster_versions.ids", "clusters.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_available_cluster_versions.ids", "clusters.0.cluster_id"),
				),
			},
		},
	})
}

const testAccKubernetesAvailableClusterVersionsDataSource_basic = testAccTkeCluster + `
data "tencentcloud_kubernetes_available_cluster_versions" "id" {
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
}
`

const testAccKubernetesAvailableClusterVersionsDataSource_multiple = testAccTkeCluster + `
data "tencentcloud_kubernetes_available_cluster_versions" "ids" {
  cluster_ids = [tencentcloud_kubernetes_cluster.managed_cluster.id]
}
`
