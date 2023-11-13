package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRabbitmqNodeListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqNodeListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rabbitmq_node_list.rabbitmq_node_list")),
			},
		},
	})
}

const testAccTdmqRabbitmqNodeListDataSource = `

data "tencentcloud_tdmq_rabbitmq_node_list" "rabbitmq_node_list" {
  instance_id = ""
  node_name = ""
  filters {
		name = ""
		values = 

  }
  sort_element = ""
  sort_order = ""
  }

`
