package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisUpgradeVersionOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisUpgradeVersionOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_upgrade_version_operation.upgrade_version_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_upgrade_version_operation.upgrade_version_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisUpgradeVersionOperation = `

resource "tencentcloud_redis_upgrade_version_operation" "upgrade_version_operation" {
  instance_id = "crs-c1nl9rpv"
  target_instance_type = "6"
  switch_option = 2
}

`
