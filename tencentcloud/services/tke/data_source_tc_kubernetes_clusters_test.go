package tke_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesClusterDataSource(t *testing.T) {
	t.Parallel()

	key := "data.tencentcloud_kubernetes_clusters.name"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTencentCloudTkeStr,
				Check: resource.ComposeTestCheckFunc(
					// name filter
					tcacctest.AccCheckTencentCloudDataSourceID(key),
					resource.TestCheckResourceAttr(key, "cluster_name", "terraform"),
					resource.TestCheckResourceAttrSet(key, "list.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudKubernetesClusterTagsDataSource(t *testing.T) {
	t.Parallel()

	key := "data.tencentcloud_kubernetes_clusters.tags"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTencentCloudTkeTags,
				Check: resource.ComposeTestCheckFunc(
					// tags filter
					tcacctest.AccCheckTencentCloudDataSourceID(key),
					resource.TestCheckResourceAttrSet(key, "list.#"),
				),
			},
		},
	})
}

const testAccDataSourceTencentCloudTkeStr = `
data "tencentcloud_kubernetes_clusters" "name" {
  #examples have been created to serve other resources
  cluster_name = "terraform"

  tags = {
    "test" = "test"
  }
}
`

const testAccDataSourceTencentCloudTkeTags = `
data "tencentcloud_kubernetes_clusters" "tags" {
  #examples have been created to serve other resources
  tags = {
    "test" = "test"
  }
}
`
