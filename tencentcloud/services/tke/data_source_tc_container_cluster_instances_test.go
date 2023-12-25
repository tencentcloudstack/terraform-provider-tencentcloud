package tke_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDataSourceContainerClusterInstances(t *testing.T) {
	t.Parallel()
	key := "data.tencentcloud_container_cluster_instances.foo_instance"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceContainerClusterInstancesConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(key),
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
	cluster_id = data.tencentcloud_container_clusters.foo.clusters.0.cluster_id
}
`
