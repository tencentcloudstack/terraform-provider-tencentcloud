package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixTdmqRabbitmqNodeListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqNodeListDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rabbitmq_node_list.rabbitmq_node_list"),
				),
			},
		},
	})
}

const testAccTdmqRabbitmqNodeListDataSource = `
data "tencentcloud_tdmq_rabbitmq_node_list" "rabbitmq_node_list" {
  instance_id = "amqp-testtesttest"
  node_name   = "keep-node"
  filters {
    name   = "nodeStatus"
    values = ["running", "down"]
  }
  sort_element = "cpuUsage"
  sort_order   = "descend"
}
`
