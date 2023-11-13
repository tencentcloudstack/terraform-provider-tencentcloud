package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisInstanceShardsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceShardsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_instance_shards.instance_shards")),
			},
		},
	})
}

const testAccRedisInstanceShardsDataSource = `

data "tencentcloud_redis_instance_shards" "instance_shards" {
  instance_id = "crs-c1nl9rpv"
  filter_slave = &lt;nil&gt;
  }

`
