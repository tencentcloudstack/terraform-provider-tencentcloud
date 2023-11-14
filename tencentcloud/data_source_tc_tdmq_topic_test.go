package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqTopicDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqTopicDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_topic.topic")),
			},
		},
	})
}

const testAccTdmqTopicDataSource = `

data "tencentcloud_tdmq_topic" "topic" {
  offset = 0
  limit = 20
  topic_name = "topic_name"
  topic_name_list = &lt;nil&gt;
  is_tag_filter = &lt;nil&gt;
  filters {
		name = "tag"
		values = &lt;nil&gt;

  }
  topic_list {
		topic_id = &lt;nil&gt;
		topic_name = &lt;nil&gt;
		msg_retention_seconds = &lt;nil&gt;
		max_msg_size = &lt;nil&gt;
		qps = &lt;nil&gt;
		filter_type = &lt;nil&gt;
		create_time = &lt;nil&gt;
		last_modify_time = &lt;nil&gt;
		msg_count = &lt;nil&gt;
		create_uin = &lt;nil&gt;
		tags {
			tag_key = &lt;nil&gt;
			tag_value = &lt;nil&gt;
		}
		trace = &lt;nil&gt;
		tenant_id = &lt;nil&gt;
		namespace_name = &lt;nil&gt;
		status = &lt;nil&gt;
		broker_type = &lt;nil&gt;

  }
  total_count = &lt;nil&gt;
  request_id = &lt;nil&gt;
}

`
