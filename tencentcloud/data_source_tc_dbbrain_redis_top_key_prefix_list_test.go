package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainRedisTopKeyPrefixListDataSource_basic(t *testing.T) {
	t.Parallel()
	// loc, _ := time.LoadLocation("Asia/Chongqing")
	// queryDate := time.Now().AddDate(0, 0, -20).In(loc).Format("2006-01-02")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainRedisTopKeyPrefixListDataSource, "2023-06-13"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list", "date", "2023-06-13"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list", "product", "redis"),
					// return
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list", "items.0.ave_element_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list", "items.0.length"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list", "items.0.key_pre_index"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list", "items.0.item_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list", "items.0.count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_key_prefix_list.redis_top_key_prefix_list", "items.0.max_element_size"),
				),
			},
		},
	})
}

const testAccDbbrainRedisTopKeyPrefixListDataSource = CommonPresetRedis + `

data "tencentcloud_dbbrain_redis_top_key_prefix_list" "redis_top_key_prefix_list" {
	instance_id = local.redis_id
	date        = "%s"
	product     = "redis"
}

`
