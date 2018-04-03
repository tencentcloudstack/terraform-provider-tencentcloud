package tencentcloud

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudDataSourceContainerClusters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceContainerClustersConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_container_clusters.foo"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceContainerClustersConfig_basic = `
data "tencentcloud_container_clusters" "foo" {
}
`

