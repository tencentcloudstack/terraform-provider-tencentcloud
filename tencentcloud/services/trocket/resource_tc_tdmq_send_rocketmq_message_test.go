package trocket_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixTdmqSendRocketmqMessageResource_basic -v
func TestAccTencentCloudNeedFixTdmqSendRocketmqMessageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqSendRocketmqMessage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_send_rocketmq_message.send_rocketmq_message", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_send_rocketmq_message.send_rocketmq_message", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_send_rocketmq_message.send_rocketmq_message", "namespace_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_send_rocketmq_message.send_rocketmq_message", "topic_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_send_rocketmq_message.send_rocketmq_message", "msg_body"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_send_rocketmq_message.send_rocketmq_message", "msg_key"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_send_rocketmq_message.send_rocketmq_message", "msg_tag"),
				),
			},
		},
	})
}

const testAccTdmqSendRocketmqMessage = `
resource "tencentcloud_tdmq_send_rocketmq_message" "send_rocketmq_message" {
  cluster_id   = "rocketmq-7k45z9dkpnne"
  namespace_id = "test_ns"
  topic_name   = "test_topic"
  msg_body     = "msg key"
  msg_key      = "msg tag"
  msg_tag      = "msg value"
}
`
