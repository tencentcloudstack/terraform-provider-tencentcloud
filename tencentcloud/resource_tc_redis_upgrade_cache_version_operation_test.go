package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisUpgradeCacheVersionOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisUpgradeCacheVersionOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_upgrade_cache_version_operation.upgrade_cache_version_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_upgrade_cache_version_operation.upgrade_cache_version_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisUpgradeCacheVersionOperation = `

resource "tencentcloud_redis_upgrade_cache_version_operation" "upgrade_cache_version_operation" {
  instance_id = "crs-c1nl9rpv"
  current_redis_version = "5.0.0"
  upgrade_redis_version = "5.0.0"
  instance_type_upgrade_now = 1
}

`
