package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaTopicSubscribeGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaTopicSubscribeGroupDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_topic_subscribe_group.topic_subscribe_group")),
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
