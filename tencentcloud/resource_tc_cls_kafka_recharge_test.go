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

resource "tencentcloud_cls_kafka_recharge" "kafka_recharge" {
  topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  name = "test"
  kafka_type = 0
  user_kafka_topics = "topic1,topic2"
  kafka_instance = "CKafka-xxxxxx"
  server_addr = "test.cls.tencentyun.com:9095"
  is_encryption_addr = false
  protocol {
		protocol = "sasl_plaintext"
		mechanism = "PLAIN"
		user_name = "username"
		password = "xxxxxx"

  }
  consumer_group_name = "group1"
  log_recharge_rule {
		recharge_type = "json_log"
		encoding_format = 0
		default_time_switch = true
		log_regex = "*"
		un_match_log_switch = true
		un_match_log_key = "test"
		un_match_log_time_src = 0
		default_time_src = 0
		time_key = "time"
		time_regex = "*"
		time_format = "%m/%d/%Y"
		time_zone = "null"
		metadata = 
		keys = 

  }
}

`
