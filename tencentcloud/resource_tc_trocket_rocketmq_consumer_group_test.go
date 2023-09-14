package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTrocketRocketmqConsumerGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTrocketRocketmqConsumerGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "consumer_group", "test_consumer_group"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "max_retry_times", "20"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "consume_enable", "false"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "consume_message_orderly", "true"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "remark", "test for terraform"),
				),
			},
			{
				Config: testAccTrocketRocketmqConsumerGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "consumer_group", "test_consumer_group"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "max_retry_times", "24"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "consume_enable", "true"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "consume_message_orderly", "true"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group", "remark", "test terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_trocket_rocketmq_consumer_group.rocketmq_consumer_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTrocketRocketmqConsumerGroup = `

resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test_consumer_group"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-3a9fo1k9"
  subnet_id     = "subnet-8nby1yxg"
  tags          = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_consumer_group" "rocketmq_consumer_group" {
  instance_id             = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  consumer_group          = "test_consumer_group"
  max_retry_times         = 20
  consume_enable          = false
  consume_message_orderly = true
  remark                  = "test for terraform"
}

`

const testAccTrocketRocketmqConsumerGroupUpdate = `

resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test_consumer_group"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-3a9fo1k9"
  subnet_id     = "subnet-8nby1yxg"
  tags          = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_consumer_group" "rocketmq_consumer_group" {
  instance_id             = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  consumer_group          = "test_consumer_group"
  max_retry_times         = 24
  consume_enable          = true
  consume_message_orderly = true
  remark                  = "test terraform"
}

`
