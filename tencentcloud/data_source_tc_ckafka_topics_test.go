package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCkafkaTopicDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudKafkaTopicDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudCkafkaTopicDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaTopicInstanceExists("tencentcloud_ckafka_topic.kafka_topic"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_id", "ckafka-f9ife4zz"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.topic_name", "ckafkaTopic-tf-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.partition_num", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.replica_num", "2"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.create_time"),
				),
			},
		},
	})
}

const testAccTencentCloudCkafkaTopicDataSourceConfig = `
resource "tencentcloud_ckafka_topic" "kafka_topic" {
   instance_id                   = "ckafka-f9ife4zz"
   topic_name                = "ckafkaTopic-tf-test"
   replica_num                   = 2
   partition_num              = 1
}

data "tencentcloud_ckafka_topics" "kafka_topics" {
	instance_id						= tencentcloud_ckafka_topic.kafka_topic.instance_id
	topic_name						= tencentcloud_ckafka_topic.kafka_topic.topic_name
}
`
