package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqDeadLetterSourceQueueDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqDeadLetterSourceQueueDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_dead_letter_source_queue.dead_letter_source_queue")),
			},
		},
	})
}

const testAccTdmqDeadLetterSourceQueueDataSource = `

data "tencentcloud_tdmq_dead_letter_source_queue" "dead_letter_source_queue" {
  dead_letter_queue_name = ""
  source_queue_name = ""
  }

`
