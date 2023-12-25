package trocket_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixTdmqRocketmqMessagesDataSource_basic -v
func TestAccTencentCloudNeedFixTdmqRocketmqMessagesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqMessageDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rocketmq_messages.message"),
				),
			},
		},
	})
}

const testAccTdmqMessageDataSource = `
data "tencentcloud_tdmq_rocketmq_messages" "message" {
  cluster_id     = "rocketmq-rkrbm52djmro"
  environment_id = "keep_ns"
  topic_name     = "keep-topic"
  msg_id         = "A9FE8D0567FE15DB97425FC08EEF0000"
  query_dlq_msg  = false
}
`
