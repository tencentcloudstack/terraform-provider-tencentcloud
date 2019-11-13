package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudDataSourceContainerClusters(t *testing.T) {
	key := "data.tencentcloud_container_clusters.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceContainerClustersConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(key),
					resource.TestCheckResourceAttrSet(key, "total_count"),
					resource.TestCheckResourceAttrSet(key, "clusters.#"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceContainerClustersConfig_basic = `
data "tencentcloud_container_clusters" "foo" {
}
`
