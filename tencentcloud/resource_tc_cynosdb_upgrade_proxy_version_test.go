package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbUpgradeProxyVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbUpgradeProxyVersion,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_upgrade_proxy_version.upgrade_proxy_version", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_upgrade_proxy_version.upgrade_proxy_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbUpgradeProxyVersion = `

resource "tencentcloud_cynosdb_upgrade_proxy_version" "upgrade_proxy_version" {
  cluster_id = "cynosdbmysql-xxxxxx"
  src_proxy_version = "无"
  dst_proxy_version = "无"
  proxy_group_id = "无"
  is_in_maintain_period = "no"
}

`
