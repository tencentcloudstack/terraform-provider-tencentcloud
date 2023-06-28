package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlProxyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlProxy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_proxy.proxy", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_proxy.proxy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlProxy = `

resource "tencentcloud_mysql_proxy" "proxy" {
	instance_id    = "cdb-fitq5t9h"
	uniq_vpc_id    = "vpc-4owdpnwr"
	uniq_subnet_id = "subnet-ahv6swf2"
	proxy_node_custom {
	  node_count = 1
	  cpu        = 2
	  mem        = 4000
	  region     = "ap-guangzhou"
	  zone       = "ap-guangzhou-3"
	}
	security_group        = ["sg-edmur627"]
	desc                  = "desc1"
	connection_pool_limit = 2
	vip                   = "172.16.17.101"
	vport                 = 3306
}

`
