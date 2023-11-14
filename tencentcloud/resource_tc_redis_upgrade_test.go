package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisUpgradeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisUpgrade,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_upgrade.upgrade", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_upgrade.upgrade",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisUpgrade = `

resource "tencentcloud_redis_upgrade" "upgrade" {
  instance_id = "crs-c1nl9rpv"
  start_time = "17:00"
  end_time = "19:00"
}

`
