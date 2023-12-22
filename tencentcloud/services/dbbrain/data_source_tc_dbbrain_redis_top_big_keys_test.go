package dbbrain_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainRedisTopBigKeysDataSource_basic(t *testing.T) {
	t.Parallel()
	// loc, _ := time.LoadLocation("Asia/Chongqing")
	// queryDate := time.Now().AddDate(0, 0, -20).In(loc).Format("2006-01-02")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainRedisTopBigKeysDataSource, "2023-06-13"),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "date", "2023-06-13"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "product", "redis"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "sort_by", "Capacity"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "key_type", "string"),
					// return
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "top_keys.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "top_keys.0.key"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "top_keys.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "top_keys.0.encoding"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "top_keys.0.expire_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "top_keys.0.length"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "top_keys.0.item_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_redis_top_big_keys.redis_top_big_keys", "top_keys.0.max_element_size"),
				),
			},
		},
	})
}

const testAccDbbrainRedisTopBigKeysDataSource = tcacctest.CommonPresetRedis + `

data "tencentcloud_dbbrain_redis_top_big_keys" "redis_top_big_keys" {
	instance_id = local.redis_id
	date        = "%s"
	product     = "redis"
	sort_by     = "Capacity"
	key_type    = "string"
}

`
