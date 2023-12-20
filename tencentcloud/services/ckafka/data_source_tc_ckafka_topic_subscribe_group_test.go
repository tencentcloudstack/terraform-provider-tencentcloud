package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaTopicSubscribeGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaTopicSubscribeGroupDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_topic_subscribe_group.topic_subscribe_group")),
			},
		},
	})
}

const testAccCkafkaTopicSubscribeGroupDataSource = `
data "tencentcloud_ckafka_topic_subscribe_group" "topic_subscribe_group" {
	instance_id = "ckafka-vv7wpvae"
	topic_name = "keep-topic"
}
`
