package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbProxyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbProxy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_proxy.proxy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbProxy = `

resource "tencentcloud_cynosdb_proxy" "proxy" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  cpu = 2
  mem = 4000
  unique_vpc_id = "无"
  unique_subnet_id = "无"
  proxy_count = 2
  connection_pool_type = "SessionConnectionPool"
  open_connection_pool = "yes"
  connection_pool_time_out = 0
  security_group_ids = 
  description = "无"
  proxy_zones {
		proxy_node_zone = ""
		proxy_node_count = 

  }
}

`
