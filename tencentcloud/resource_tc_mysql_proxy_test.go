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
  instance_id = ""
  uniq_vpc_id = ""
  uniq_subnet_id = ""
  proxy_node_custom {
		node_count = 
		cpu = 
		mem = 
		region = ""
		zone = ""

  }
  security_group = 
  desc = ""
  connection_pool_limit = 
}

`
