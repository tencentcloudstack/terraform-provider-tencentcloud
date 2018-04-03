package tencentcloud

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudDataSourceContainerClusterInstances(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceContainerClusterInstancesConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_container_cluster_instances.foo_instance"),
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

