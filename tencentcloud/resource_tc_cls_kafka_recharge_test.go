package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClsKafkaRechargeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsKafkaRecharge,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cls_kafka_recharge.kafka_recharge", "id")),
			},
			{
				ResourceName:      "tencentcloud_cls_kafka_recharge.kafka_recharge",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
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
