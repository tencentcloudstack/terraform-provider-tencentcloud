package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqQueueResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqQueue,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_queue.queue", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_queue.queue",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqQueue = `

resource "tencentcloud_tdmq_queue" "queue" {
  queue_name = "queue_name"
  max_msg_heap_num = 10000000
  polling_wait_seconds = 0
  visibility_timeout = 30
  max_msg_size = 65536
  msg_retention_seconds = 3600
  rewind_seconds = 0
  transaction = 1
  first_query_interval = 1
  max_query_count = 1
  dead_letter_queue_name = "dead_letter_queue_name"
  policy = 0
  max_receive_count = 50
  max_time_to_live = 300
  trace = false
  retention_size_in_m_b = 0
}

`
