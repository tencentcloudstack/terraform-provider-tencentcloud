package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCatNodeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCatNodeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cat_node.node")),
			},
		},
	})
}

const testAccCatNodeDataSource = `

data "tencentcloud_cat_node" "node" {
  node_type = 
  location = 
  is_i_pv6 = 
  node_name = ""
  pay_mode = 
  task_type = 
  }

`
