package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainRedisTopKeyPrefixListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainRedisTopKeyPrefixListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list")),
			},
		},
	})
}

const testAccDbbrainRedisTopKeyPrefixListDataSource = `

data "tencentcloud_dbbrain_redis_top_key_prefix_list" "redis_top_key_prefix_list" {
  instance_id = ""
  date = ""
  product = ""
    }

`
