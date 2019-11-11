package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudTke(t *testing.T) {

	key := "data.tencentcloud_kubernetes_clusters.name"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTencentCloudTkeStr,
				Check: resource.ComposeTestCheckFunc(
					// name filter
					testAccCheckTencentCloudDataSourceID(key),
					resource.TestCheckResourceAttr(key, "cluster_name", "terraform"),
					resource.TestCheckResourceAttrSet(key, "list.#"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudTkeTags(t *testing.T) {

	key := "data.tencentcloud_kubernetes_clusters.tags"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTencentCloudTkeTags,
				Check: resource.ComposeTestCheckFunc(
					// tags filter
					testAccCheckTencentCloudDataSourceID(key),
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
