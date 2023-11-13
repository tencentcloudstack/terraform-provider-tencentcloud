package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbProxyEndPointResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbProxyEndPoint,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy_end_point.proxy_end_point", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_proxy_end_point.proxy_end_point",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbProxyEndPoint = `

resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  unique_vpc_id = "无"
  unique_subnet_id = "无"
  connection_pool_type = "SessionConnectionPool"
  open_connection_pool = "yes"
  connection_pool_time_out = 0
  security_group_ids = 
  description = "无"
  vip = "无"
  weight_mode = "无"
  auto_add_ro = "无"
  fail_over = "无"
  consistency_type = "无"
  rw_type = "无"
  consistency_time_out = 无
  trans_split = 无
  access_mode = "无"
  instance_weights {
		instance_id = ""
		weight = 

  }
}

`
