package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCynosdbUpgradeProxyVersionResource_basic -v
func TestAccTencentCloudNeedFixCynosdbUpgradeProxyVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbUpgradeProxyVersion,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_upgrade_proxy_version.upgrade_proxy_version", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_upgrade_proxy_version.upgrade_proxy_version", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_upgrade_proxy_version.upgrade_proxy_version", "dst_proxy_version"),
				),
			},
		},
	})
}

const testAccCynosdbUpgradeProxyVersion = `
resource "tencentcloud_cynosdb_upgrade_proxy_version" "upgrade_proxy_version" {
  cluster_id = "cynosdbmysql-bws8h88b"
  dst_proxy_version = "1.3.7"
}
`
