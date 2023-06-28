package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlUpgradeProxyVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlUpgradeProxyVersion,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_upgrade_proxy_version.upgrade_proxy_version", "id"),
				),
			},
		},
	})
}

const testAccMysqlUpgradeProxyVersion = `

resource "tencentcloud_mysql_upgrade_proxy_version" "upgrade_proxy_version" {
	instance_id       = "cdb-fitq5t9h"
	proxy_group_id    = "proxy-h1ub486b"
	src_proxy_version = "1.3.6"
	dst_proxy_version = "1.3.7"
	upgrade_time      = "nowTime"
}

`
