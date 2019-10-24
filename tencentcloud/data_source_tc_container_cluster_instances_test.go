package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudDataSourceContainerClusterInstances(t *testing.T) {
	key := "data.tencentcloud_container_cluster_instances.foo_instance"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceContainerClusterInstancesConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(key),
					resource.TestCheckResourceAttrSet(key, "total_count"),
					resource.TestCheckResourceAttrSet(key, "nodes.#"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceContainerClusterInstancesConfig_basic = `
data "tencentcloud_container_clusters" "foo" {
}

data "tencentcloud_container_cluster_instances" "foo_instance" {
	cluster_id = "${data.tencentcloud_container_clusters.foo.clusters.0.cluster_id}"
}
`
