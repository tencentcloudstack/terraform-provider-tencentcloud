package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisInstanceZoneInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceZoneInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_instance_zone_info.instance_zone_info")),
			},
		},
	})
}

const testAccRedisInstanceZoneInfoDataSource = `

data "tencentcloud_redis_instance_zone_info" "instance_zone_info" {
  instance_id = "crs-c1nl9rpv"
  }

`
