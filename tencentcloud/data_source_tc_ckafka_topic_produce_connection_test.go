package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaTopicProduceConnectionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaTopicProduceConnectionDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ckafka_topic_produce_connection.topic_produce_connection")),
			},
		},
	})
}

const testAccCkafkaTopicProduceConnectionDataSource = `

data "tencentcloud_ckafka_topic_produce_connection" "topic_produce_connection" {
  instance_id = "InstanceId"
  topic_name = "TopicName"
  }

`
