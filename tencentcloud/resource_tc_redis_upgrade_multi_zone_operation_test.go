package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisUpgradeMultiZoneOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisUpgradeMultiZoneOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_upgrade_multi_zone_operation.upgrade_multi_zone_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_upgrade_multi_zone_operation.upgrade_multi_zone_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisUpgradeMultiZoneOperation = `

resource "tencentcloud_redis_upgrade_multi_zone_operation" "upgrade_multi_zone_operation" {
  instance_id = "crs-c1nl9rpv"
  upgrade_proxy_and_redis_server = 
}

`
