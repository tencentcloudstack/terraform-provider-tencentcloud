package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqQueueDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqQueueDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_queue.queue")),
			},
		},
	})
}

const testAccTdmqQueueDataSource = `

data "tencentcloud_tdmq_queue" "queue" {
  offset = 0
  limit = 20
  queue_name = "queue_name"
  queue_name_list = &lt;nil&gt;
  is_tag_filter = true
  filters {
		name = "tag"
		values = &lt;nil&gt;

  }
  total_count = &lt;nil&gt;
  queue_list {
		queue_id = &lt;nil&gt;
		queue_name = &lt;nil&gt;
		qps = &lt;nil&gt;
		bps = &lt;nil&gt;
		max_delay_seconds = &lt;nil&gt;
		max_msg_heap_num = &lt;nil&gt;
		polling_wait_seconds = &lt;nil&gt;
		msg_retention_seconds = &lt;nil&gt;
		visibility_timeout = &lt;nil&gt;
		max_msg_size = &lt;nil&gt;
		rewind_seconds = &lt;nil&gt;
		create_time = &lt;nil&gt;
		last_modify_time = &lt;nil&gt;
		active_msg_num = &lt;nil&gt;
		inactive_msg_num = &lt;nil&gt;
		delay_msg_num = &lt;nil&gt;
		rewind_msg_num = &lt;nil&gt;
		min_msg_time = &lt;nil&gt;
		transaction = &lt;nil&gt;
		dead_letter_source {
			queue_id = &lt;nil&gt;
			queue_name = &lt;nil&gt;
		}
		dead_letter_policy {
			dead_letter_queue = &lt;nil&gt;
			policy = &lt;nil&gt;
			max_time_to_live = &lt;nil&gt;
			max_receive_count = &lt;nil&gt;
		}
		transaction_policy {
			first_query_interval = &lt;nil&gt;
			max_query_count = &lt;nil&gt;
		}
		create_uin = &lt;nil&gt;
		tags {
			tag_key = &lt;nil&gt;
			tag_value = &lt;nil&gt;
		}
		trace = &lt;nil&gt;
		tenant_id = &lt;nil&gt;
		namespace_name = &lt;nil&gt;
		status = &lt;nil&gt;
		max_unacked_msg_num = &lt;nil&gt;
		max_msg_backlog_size = &lt;nil&gt;
		retention_size_in_m_b = &lt;nil&gt;

  }
  request_id = &lt;nil&gt;
}

`
