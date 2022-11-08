package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTdmqRabbitmqQueue_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqQueue,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_queue.rabbitmq_queue", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_queue.rabbitmqQueue",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRabbitmqQueue = `

resource "tencentcloud_tdmq_rabbitmq_queue" "rabbitmq_queue" {
  queue = ""
  cluster_id = ""
  v_host_id = ""
  auto_delete = ""
  remark = ""
  dead_letter_exchange = ""
  dead_letter_routing_key = ""
}

`
