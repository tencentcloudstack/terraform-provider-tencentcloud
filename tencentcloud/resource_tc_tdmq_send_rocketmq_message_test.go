package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqSendRocketmqMessageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqSendRocketmqMessage,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_send_rocketmq_message.send_rocketmq_message", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_send_rocketmq_message.send_rocketmq_message",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqSendRocketmqMessage = `

resource "tencentcloud_tdmq_send_rocketmq_message" "send_rocketmq_message" {
  cluster_id = ""
  namespace_id = ""
  topic_name = ""
  msg_body = ""
  msg_key = ""
  msg_tag = ""
}

`
