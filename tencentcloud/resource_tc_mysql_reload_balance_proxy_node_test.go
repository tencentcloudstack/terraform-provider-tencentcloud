package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlReloadBalanceProxyNodeResource_basic -v
func TestAccTencentCloudMysqlReloadBalanceProxyNodeResource_basic(t *testing.T) {
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

const testAccMysqlReloadBalanceProxyNodeVar = testAccMysqlProxyVar + `

resource "tencentcloud_mysql_proxy" "proxy" {
	instance_id    = var.instance_id
	uniq_vpc_id    = var.vpc_id
	uniq_subnet_id = var.subnet_id
	proxy_node_custom {
	  node_count = 1
	  cpu        = 2
	  mem        = 4000
	  region     = "ap-guangzhou"
	  zone       = "ap-guangzhou-3"
	}
	security_group        = ["sg-edmur627"]
	desc                  = "desc"
	connection_pool_limit = 1
}

`

const testAccMysqlReloadBalanceProxyNode = testAccMysqlReloadBalanceProxyNodeVar + `

resource "tencentcloud_mysql_reload_balance_proxy_node" "reload_balance_proxy_node" {
  proxy_group_id = tencentcloud_mysql_proxy.proxy.proxy_group_id
  proxy_address_id = tencentcloud_mysql_proxy.proxy.proxy_address_id
}

`
