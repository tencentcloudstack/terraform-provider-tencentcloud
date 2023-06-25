package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlReloadBalanceProxyNodeResource_basic -v
func TestAccTencentCloudMysqlReloadBalanceProxyNodeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlReloadBalanceProxyNode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_reload_balance_proxy_node.reload_balance_proxy_node", "id"),
				),
			},
		},
	})
}

const testAccMysqlReloadBalanceProxyNode = `

resource "tencentcloud_mysql_reload_balance_proxy_node" "reload_balance_proxy_node" {
  proxy_group_id = "proxy-gmi1f78l"
  proxy_address_id = "proxyaddr-4wc4y1pq"
}

`
