package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseZookeeperReplicasDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseZookeeperReplicasDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_zookeeper_replicas.zookeeper_replicas")),
			},
		},
	})
}

const testAccTseZookeeperReplicasDataSource = `

data "tencentcloud_tse_zookeeper_replicas" "zookeeper_replicas" {
  instance_id = "ins-xxxxxx"
    tags = {
    "createdBy" = "terraform"
  }
}

`
