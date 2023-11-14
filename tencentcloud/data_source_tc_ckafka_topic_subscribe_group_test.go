package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaTopicSubscribeGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
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
  instance_id = "InstanceId"
  topic_name = "TopicName"
  }

`
