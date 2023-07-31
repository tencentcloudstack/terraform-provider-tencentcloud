package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesAvailableClusterVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKubernetesAvailableClusterVersionsDataSource_basic, defaultTkeClusterId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_available_cluster_versions.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_available_cluster_versions.id", "versions.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_available_cluster_versions.id", "clusters.#"),
				),
			},
			{
				Config: fmt.Sprintf(testAccKubernetesAvailableClusterVersionsDataSource_multiple, defaultTkeClusterId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_available_cluster_versions.ids"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_available_cluster_versions.ids", "clusters.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_available_cluster_versions.ids", "clusters.0.cluster_id"),
				),
			},
		},
	})
}

const testAccKubernetesAvailableClusterVersionsDataSource_basic = `
variable "env_default_tke_cluster_id" {
  type = string
}

data "tencentcloud_kubernetes_available_cluster_versions" "id" {
  cluster_id = var.env_default_tke_cluster_id != "" ? var.env_default_tke_cluster_id : "%s"
}

output "versions"{
  value = data.tencentcloud_kubernetes_available_cluster_versions.id.versions
}

`

const testAccKubernetesAvailableClusterVersionsDataSource_multiple = `
variable "env_default_tke_cluster_id" {
	type = string
  }

data "tencentcloud_kubernetes_available_cluster_versions" "ids" {
  cluster_ids = [var.env_default_tke_cluster_id != "" ? var.env_default_tke_cluster_id : "%s"]
}

`
