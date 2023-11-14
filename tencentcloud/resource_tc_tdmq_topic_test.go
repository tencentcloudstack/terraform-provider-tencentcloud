package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqTopicResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqTopic,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_topic.topic", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_topic.topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqTopic = `

resource "tencentcloud_tdmq_topic" "topic" {
  topic_name = "topic_name"
  max_msg_size = 65536
  filter_type = 1
  msg_retention_seconds = 86400
  trace = true
}

`
