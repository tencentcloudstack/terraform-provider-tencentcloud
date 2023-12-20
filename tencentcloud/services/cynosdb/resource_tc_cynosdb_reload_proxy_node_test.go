package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbReloadProxyNodeResource_basic -v
func TestAccTencentCloudCynosdbReloadProxyNodeResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckCynosdbProxyDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbReloadProxyNode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_reload_proxy_node.reload_proxy_node", "id")),
			},
		},
	})
}

const testAccCynosdbReloadProxyNode = testAccCynosdbProxy + `
resource "tencentcloud_cynosdb_reload_proxy_node" "reload_proxy_node" {
  cluster_id     = tencentcloud_cynosdb_proxy.proxy.cluster_id
  proxy_group_id = tencentcloud_cynosdb_proxy.proxy.proxy_group_id
}
`
