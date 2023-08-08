package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccTencentCloudClsKafkaRechargeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckClsKafkaRechargeDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsKafkaRecharge,
				Check: resource.ComposeTestCheckFunc(testAccCheckClsKafkaRechargeExists("tencentcloud_cls_kafka_recharge.kafka_recharge"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_kafka_recharge.kafka_recharge", "topic_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_kafka_recharge.kafka_recharge", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_kafka_recharge.kafka_recharge", "kafka_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_kafka_recharge.kafka_recharge", "offset"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_kafka_recharge.kafka_recharge", "user_kafka_topics"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_kafka_recharge.kafka_recharge", "kafka_instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_kafka_recharge.kafka_recharge", "log_recharge_rule.#")),
			},
			{
				ResourceName:      "tencentcloud_cls_kafka_recharge.kafka_recharge",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClsKafkaRechargeDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clsService := ClsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cls_kafka_recharge" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		kafkaRechargeId := idSplit[0]
		kafkaTopic := idSplit[1]

		instance, err := clsService.DescribeClsKafkaRechargeById(ctx, kafkaRechargeId, kafkaTopic)
		if err != nil {
			continue
		}
		if instance != nil {
			return fmt.Errorf("[CHECK][CLS KafkaRecharge][Destroy] check: CLS KafkaRecharge still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClsKafkaRechargeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS KafkaRecharge][Exists] check: CLS KafkaRecharge %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS dataTransform][Create] check: CLS KafkaRecharge id is not set")
		}
		clsService := ClsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)

		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		kafkaRechargeId := idSplit[0]
		kafkaTopic := idSplit[1]

		instance, err := clsService.DescribeClsKafkaRechargeById(ctx, kafkaRechargeId, kafkaTopic)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CLS KafkaRecharge][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsKafkaRecharge = `
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-example-logset1"
  tags = {
    "createdBy" = "terraform"
  }
}
resource "tencentcloud_cls_topic" "topic" {
  topic_name           = "tf-example-topic"
  logset_id            = tencentcloud_cls_logset.logset.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test",
  }
}

resource "tencentcloud_cls_kafka_recharge" "kafka_recharge" {
  topic_id = tencentcloud_cls_topic.topic.id
  name = "tf-example-recharge"
  kafka_type = 0
  offset = -2
  user_kafka_topics = "dasdasd"
  kafka_instance = "ckafka-qzoeaqx8"
  log_recharge_rule {
    recharge_type = "json_log"
    encoding_format = 0
    default_time_switch = true
  }
}
`
