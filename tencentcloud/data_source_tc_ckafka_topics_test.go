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
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.topic_name", "ckafkaTopic-tf-test06"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.partition_num", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.replica_num", "2"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.note"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.create_time"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.enable_white_list"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.ip_white_list_count"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.forward_interval"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.forward_status"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.config.0.retention"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.config.0.sync_replica_min_num"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.config.0.clean_up_policy"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.config.0.unclean_leader_election_enable"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.config.0.max_message_bytes"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.config.0.segment"),
					//resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.config.0.segment_bytes"),
				),
			},
		},
	})
}

const testAccTencentCloudCkafkaTopicDataSourceConfig = `

resource "tencentcloud_ckafka_topic" "kafka_topic" {
   instance_id                   = "ckafka-f9ife4zz"
   topic_name                = "ckafkaTopic-tf-test06"
   replica_num                   = 2
   partition_num              = 1
}

data "tencentcloud_ckafka_topics" "kafka_topics" {
	instance_id						= tencentcloud_ckafka_topic.kafka_topic.instance_id
	topic_name						= tencentcloud_ckafka_topic.kafka_topic.topic_name
}
`
