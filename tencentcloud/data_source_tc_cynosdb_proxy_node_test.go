package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbProxyNodeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbProxyNodeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_proxy_node.proxy_node")),
			},
		},
	})
}

const testAccCynosdbProxyNodeDataSource = `

data "tencentcloud_cynosdb_proxy_node" "proxy_node" {
  order_by = "&quot;CREATETIME&quot;"
  order_by_type = "&quot;desc&quot;"
  filters {
		names = 
		values = 
		exact_match = 
		name = ""
		operator = ""

  }
  }

`
