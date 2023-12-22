package trabbit_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixTdmqRabbitmqNodeListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqNodeListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rabbitmq_node_list.rabbitmq_node_list"),
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
