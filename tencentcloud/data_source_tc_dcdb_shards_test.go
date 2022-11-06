package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDcdbShardsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbShards_basic, defaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_shards.shards"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_shards.shards", "list.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_shards.shards", "list.0.instance_id", defaultDcdbInstanceId),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shards.shards", "list.0.shard_instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_shards.shards", "list.1.instance_id", defaultDcdbInstanceId),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shards.shards", "list.1.shard_instance_id"),
				),
			},
		},
	})
}

const testAccDataSourceDcdbShards_basic = `
data "tencentcloud_dcdb_instances" "instances" {
	instance_ids = ["%s"]
}

data "tencentcloud_dcdb_shards" "shards" {
	instance_id = data.tencentcloud_dcdb_instances.instances.list.0.instance_id
	shard_instance_ids = [data.tencentcloud_dcdb_instances.instances.list.0.shard_detail.0.shard_instance_id, data.tencentcloud_dcdb_instances.instances.list.0.shard_detail.1.shard_instance_id]
}

`
