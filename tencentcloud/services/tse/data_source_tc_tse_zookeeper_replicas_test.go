package tse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTseZookeeperReplicasDataSource_basic -v
func TestAccTencentCloudTseZookeeperReplicasDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseZookeeperReplicasDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tse_zookeeper_replicas.zookeeper_replicas"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_zookeeper_replicas.zookeeper_replicas", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_zookeeper_replicas.zookeeper_replicas", "replicas.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_zookeeper_replicas.zookeeper_replicas", "replicas.0.alias_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_zookeeper_replicas.zookeeper_replicas", "replicas.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_zookeeper_replicas.zookeeper_replicas", "replicas.0.role"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_zookeeper_replicas.zookeeper_replicas", "replicas.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_zookeeper_replicas.zookeeper_replicas", "replicas.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_zookeeper_replicas.zookeeper_replicas", "replicas.0.zone_id"),
				),
			},
		},
	})
}

const testAccTseZookeeperReplicasDataSource = testAccTseInstance + `

data "tencentcloud_tse_zookeeper_replicas" "zookeeper_replicas" {
  instance_id = tencentcloud_tse_instance.instance.id
}

`
