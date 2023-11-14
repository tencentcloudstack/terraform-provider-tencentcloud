package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbProxyRwSplitResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbProxyRwSplit,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy_rw_split.proxy_rw_split", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_proxy_rw_split.proxy_rw_split",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbProxyRwSplit = `

resource "tencentcloud_cynosdb_proxy_rw_split" "proxy_rw_split" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  proxy_group_id = "无"
  consistency_type = "无"
  consistency_time_out = "无"
  weight_mode = "无"
  instance_weights {
		instance_id = ""
		weight = 

  }
  fail_over = "无"
  auto_add_ro = "无"
  open_rw = "无"
  rw_type = "无"
  trans_split = 无
  access_mode = "无"
  open_connection_pool = "无"
  connection_pool_type = "无"
  connection_pool_time_out = 无
}

`
