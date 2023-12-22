package dcdb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDCDBShardsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbShards_basic, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_shards.shards"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_shards.shards", "list.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_shards.shards", "list.0.instance_id", tcacctest.DefaultDcdbInstanceId),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_shards.shards", "list.0.shard_instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_shards.shards", "list.1.instance_id", tcacctest.DefaultDcdbInstanceId),
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
