package tke_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesAvailableClusterVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKubernetesAvailableClusterVersionsDataSource_basic, tcacctest.DefaultTkeClusterId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_available_cluster_versions.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_available_cluster_versions.id", "versions.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_available_cluster_versions.id", "clusters.#"),
				),
			},
			{
				Config: fmt.Sprintf(testAccKubernetesAvailableClusterVersionsDataSource_multiple, tcacctest.DefaultTkeClusterId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_available_cluster_versions.ids"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_available_cluster_versions.ids", "clusters.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_available_cluster_versions.ids", "clusters.0.cluster_id", tcacctest.DefaultTkeClusterId),
				),
			},
		},
	})
}

const testAccKubernetesAvailableClusterVersionsDataSource_basic = `

data "tencentcloud_kubernetes_available_cluster_versions" "id" {
  cluster_id = "%s"
}

output "versions"{
  value = data.tencentcloud_kubernetes_available_cluster_versions.id.versions
}

`

const testAccKubernetesAvailableClusterVersionsDataSource_multiple = `

data "tencentcloud_kubernetes_available_cluster_versions" "ids" {
  cluster_ids = ["%s"]
}

`
