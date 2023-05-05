package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataSourceTcmqQueue = "data.tencentcloud_tcmq_queue.queue"

func TestAccTencentCloudTcmqQueueDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcmqQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudTcmqQueueDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSourceTcmqQueue, "queue_list.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudTcmqQueueDataSource_basic = `
resource "tencentcloud_tcmq_queue" "queue" {
	queue_name="test_queue_datasource"
}
  
data "tencentcloud_tcmq_queue" "queue" {
	queue_name = tencentcloud_tcmq_queue.queue.queue_name
}
`
