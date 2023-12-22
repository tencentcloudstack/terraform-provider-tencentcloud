package ckafka_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudInternationalCkafkaResource_topic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalKafkaTopicInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaTopicInstanceExists("tencentcloud_ckafka_topic.kafka_topic"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_topic.kafka_topic", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "topic_name", "ckafka-topic-tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "note", "this is test ckafka topic"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "replica_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "partition_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "enable_white_list", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "ip_white_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "ip_white_list.0", "192.168.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "clean_up_policy", "delete"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "sync_replica_min_num", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_topic.kafka_topic", "unclean_leader_election_enable"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "segment", "86400000"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(60 * time.Second)
				},
				ResourceName:      "tencentcloud_ckafka_topic.kafka_topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccInternationalKafkaTopicInstance = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance_postpaid" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = 100007
  vpc_id             = var.international_vpc_id
  subnet_id          = var.international_subnet_id
  msg_retention_time = 1300
  kafka_version      = "1.1.1"
  disk_size          = 500
  band_width         = 20
  disk_type          = "CLOUD_BASIC"
  partition          = 400
  charge_type        = "POSTPAID_BY_HOUR"

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }
}

resource "tencentcloud_ckafka_topic" "kafka_topic" {
	instance_id                     = tencentcloud_ckafka_instance.kafka_instance_postpaid.id
	topic_name                          = "ckafka-topic-tf-test"
	note                                = "this is test ckafka topic"
	replica_num                         = 2
	partition_num                       = 2
	enable_white_list                   = true
	ip_white_list                       = ["192.168.1.1"]
	clean_up_policy                     = "delete"
	sync_replica_min_num                = 1
	unclean_leader_election_enable      = false
	segment                             = 86400000
	max_message_bytes                   = 1024
}
`
