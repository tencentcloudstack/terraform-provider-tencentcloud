package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaTopicDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudKafkaTopicDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudCkafkaTopicDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaTopicInstanceExists("tencentcloud_ckafka_topic.kafka_topic"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_id", "ckafka-vv7wpvae"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.topic_name", "ckafkaTopic-tf-test"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.partition_num", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.replica_num", "2"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.note", "test topic"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.enable_white_list", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.ip_white_list_count", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.clean_up_policy", "delete"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.sync_replica_min_num", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.unclean_leader_election_enable"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.segment", "86400000"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_topics.kafka_topics", "instance_list.0.retention", "60000"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.foo", "instance_list.0.partition_num"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.foo", "instance_list.0.replica_num"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_topics.foo", "instance_list.0.create_time"),
				),
			},
		},
	})
}

const testAccTencentCloudCkafkaTopicDataSourceConfig = defaultKafkaVariable + `
resource "tencentcloud_ckafka_topic" "kafka_topic" {
	instance_id                     = var.instance_id
	topic_name                      = "ckafkaTopic-tf-test"
	replica_num                     = 2
	partition_num                   = 1
	note                            = "test topic"
	enable_white_list               = true
	ip_white_list                   = ["192.168.1.1"]
	clean_up_policy                 = "delete"
	sync_replica_min_num            = 1
	unclean_leader_election_enable  = false
	segment                         = 86400000
	retention                       = 60000
	max_message_bytes               = 8388608
}

data "tencentcloud_ckafka_topics" "kafka_topics" {
	instance_id						= tencentcloud_ckafka_topic.kafka_topic.instance_id
	topic_name						= tencentcloud_ckafka_topic.kafka_topic.topic_name
}

data "tencentcloud_ckafka_topics" "foo" {
	instance_id						= tencentcloud_ckafka_topic.kafka_topic.instance_id
}
`
