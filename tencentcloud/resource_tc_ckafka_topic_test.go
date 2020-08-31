package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
)

func TestAccTencentCloudKafkaTopic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudKafkaTopicDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaTopicInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaTopicInstanceExists("tencentcloud_ckafka_topic.kafka_topic"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "instance_id", "ckafka-f9ife4zz"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "topic_name", "topic-tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "note", "this is test ckafka topic"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "replica_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "partition_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "enable_white_list", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "ip_white_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "ip_white_list.0", "192.168.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "clean_up_policy", "delete"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "sync_replica_min_num", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_topic.kafka_topic", "unclean_leader_election_enable"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "segment", "3600000"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "retention", "60000"),
				),
			},
			{
				Config: testAccKafkaTopicInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaTopicInstanceExists("tencentcloud_ckafka_topic.kafka_topic"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "instance_id", "ckafka-f9ife4zz"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "note", "this is test topic_update"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_topic.kafka_topic", "partition_num"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "enable_white_list", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "clean_up_policy", "compact"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "sync_replica_min_num", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_topic.kafka_topic", "unclean_leader_election_enable"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "segment", "4000000"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "retention", "70000"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "max_message_bytes", "8388608"),
				),
			},
		},
	})
}

func testAccTencentCloudKafkaTopicDestory(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkcService := CkafkaService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, r := range s.RootModule().Resources {
		if r.Type != "tencentcloud_ckafka_topic" {
			continue
		}

		split := strings.Split(r.Primary.ID, FILED_SP)
		if len(split) < 2 {
			continue
		}
		//check ckafka instance
		has, err := ckafkcService.DescribeCkafkaById(ctx, split[0])
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("ckafka %s is not exists", split[0])
		}
		//check ckafka topic
		var topicinfo *ckafka.TopicDetail
		outErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			topic, _, err := ckafkcService.DescribeCkafkaTopicByName(ctx, split[0], split[1])
			topicinfo = topic
			if err != nil {
				return nil
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
		if topicinfo == nil {
			return nil
		}
		if *topicinfo.TopicName == split[0] {
			return fmt.Errorf("ckafka topic %s is still existing", split[1])
		}
	}
	return nil
}

func testAccCheckKafkaTopicInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ckafka topic %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("topic instance id is not set")
		}
		ckafkcService := CkafkaService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		split := strings.Split(rs.Primary.ID, FILED_SP)
		if len(split) < 2 {
			return fmt.Errorf("ckafka topic is not set")
		}
		var topicinfo *ckafka.TopicDetail
		outErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, _, inErr := ckafkcService.DescribeCkafkaTopicByName(ctx, split[0], split[1])
			if inErr != nil {
				return retryError(inErr)

			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
		if topicinfo == nil {
			return nil
		}
		if *topicinfo.TopicName == split[0] {
			return fmt.Errorf("ckafka topic %s is still existing", split[1])
		}
		return fmt.Errorf("ckafka topic %s is not found", rs.Primary.ID)
	}
}

const testAccKafkaTopicInstance = `
resource "tencentcloud_ckafka_topic" "kafka_topic" {
	instance_id						= "ckafka-f9ife4zz"
	topic_name						= "topic-tf-test"
	note							= "this is test ckafka topic"
	replica_num						= 2
	partition_num					= 1
	enable_white_list				= 1
	ip_white_list    				= ["192.168.1.1"]
	clean_up_policy					= "delete"
	sync_replica_min_num			= 1
	unclean_leader_election_enable  = false
	segment							= 3600000
	retention						= 60000
	max_message_bytes				= 0
}
`

const testAccKafkaTopicInstanceUpdate = `
resource "tencentcloud_ckafka_topic" "kafka_topic" {
	instance_id						= "ckafka-f9ife4zz"
	topic_name						= "topic-tf-test"
	note							= "this is test topic_update"
	replica_num						= 2
	partition_num					= 1
	enable_white_list				= 1
	ip_white_list    				= ["192.168.1.2"]
	clean_up_policy					= "compact"
	sync_replica_min_num			= 2
	unclean_leader_election_enable	= true
	segment							= 4000000
	retention						= 70000
	max_message_bytes				= 8388608
}
`
