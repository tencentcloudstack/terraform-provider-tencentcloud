package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainRedisTopBigKeysDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainRedisTopBigKeysDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys")),
			},
		},
	})
}

const testAccDbbrainRedisTopBigKeysDataSource = `

data "tencentcloud_dbbrain_redis_top_big_keys" "redis_top_big_keys" {
  instance_id = ""
  date = ""
  product = ""
  sort_by = ""
  key_type = ""
    }

`
