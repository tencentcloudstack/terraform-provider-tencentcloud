package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqSubscribeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqSubscribeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_subscribe.subscribe")),
			},
		},
	})
}

const testAccTdmqSubscribeDataSource = `

data "tencentcloud_tdmq_subscribe" "subscribe" {
  topic_name = "topic_name"
  offset = 0
  limit = 20
  subscription_name = &lt;nil&gt;
  total_count = &lt;nil&gt;
  subscription_set {
		subscription_name = &lt;nil&gt;
		subscription_id = &lt;nil&gt;
		topic_owner = &lt;nil&gt;
		msg_count = &lt;nil&gt;
		last_modify_time = &lt;nil&gt;
		create_time = &lt;nil&gt;
		binding_key = &lt;nil&gt;
		endpoint = &lt;nil&gt;
		filter_tags = &lt;nil&gt;
		protocol = &lt;nil&gt;
		notify_strategy = &lt;nil&gt;
		notify_content_format = &lt;nil&gt;

  }
  request_id = &lt;nil&gt;
}

`
