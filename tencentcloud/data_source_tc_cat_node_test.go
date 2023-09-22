package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCatNodeDataSource -v
func TestAccTencentCloudCatNodeDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCatNode,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cat_node.node"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.0.city"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.0.code"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.0.code_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.0.district"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.0.ip_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.0.location"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.0.net_service"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.0.node_define_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.0.task_types.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cat_node.node", "node_define.0.type"),
				),
			},
		},
	})
}

const testAccDataSourceCatNode = `

data "tencentcloud_cat_node" "node"{
  node_type = 1
  location  = 2
  is_ipv6   = false
}

`
