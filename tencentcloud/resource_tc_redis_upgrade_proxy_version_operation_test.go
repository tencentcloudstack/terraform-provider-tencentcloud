package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixRedisUpgradeProxyVersionOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisUpgradeProxyVersionOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_upgrade_proxy_version_operation.upgrade_proxy_version_operation", "id")),
			},
		},
	})
}

const testAccRedisUpgradeProxyVersionOperation = `

resource "tencentcloud_redis_upgrade_proxy_version_operation" "upgrade_proxy_version_operation" {
  instance_id = "crs-c1nl9rpv"
  current_proxy_version = "5.0.0"
  upgrade_proxy_version = "5.0.0"
  instance_type_upgrade_now = 1
}

`
