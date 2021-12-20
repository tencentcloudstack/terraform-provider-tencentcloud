package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDataSourceContainerClusters(t *testing.T) {
	t.Parallel()
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
