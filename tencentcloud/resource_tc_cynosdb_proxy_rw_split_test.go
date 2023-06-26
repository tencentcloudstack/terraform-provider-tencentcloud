package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCynosdbProxyRwSplitResource_basic -v
func TestAccTencentCloudNeedFixCynosdbProxyRwSplitResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbProxyRwSplit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy_rw_split.proxy_rw_split", "id"),
				),
			},
		},
	})
}

const testAccCynosdbProxyRwSplit = `
resource "tencentcloud_cynosdb_proxy_rw_split" "proxy_rw_split" {
  cluster_id           = "cynosdbmysql-cgd2gpwr"
  proxy_group_id       = "cynosdbmysql-proxy-l6zf9t30"
  consistency_type     = "global"
  consistency_time_out = "30"
  weight_mode          = "system"
  instance_weights {
    instance_id = "cynosdbmysql-ins-9810be9i"
    weight      = 0
  }
  fail_over                = "yes"
  auto_add_ro              = "no"
  open_rw                  = "yes"
  rw_type                  = "READWRITE"
  trans_split              = false
  access_mode              = "balance"
  open_connection_pool     = "yes"
  connection_pool_type     = "SessionConnectionPool"
  connection_pool_time_out = 30
}
`
