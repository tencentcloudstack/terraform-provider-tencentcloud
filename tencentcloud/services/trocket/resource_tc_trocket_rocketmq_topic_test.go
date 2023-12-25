package trocket_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTrocketRocketmqTopicResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTrocketRocketmqTopic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_topic.rocketmq_topic", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_topic.rocketmq_topic", "topic_type", "NORMAL"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_topic.rocketmq_topic", "queue_num", "4"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_topic.rocketmq_topic", "remark", "test for terraform"),
				),
			},
			{
				Config: testAccTrocketRocketmqTopicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_topic.rocketmq_topic", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_topic.rocketmq_topic", "topic_type", "NORMAL"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_topic.rocketmq_topic", "queue_num", "5"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_topic.rocketmq_topic", "remark", "test terraform"),
				)},
			{
				ResourceName:      "tencentcloud_trocket_rocketmq_topic.rocketmq_topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTrocketRocketmqTopic = `

resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test_topic"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-3a9fo1k9"
  subnet_id     = "subnet-8nby1yxg"
  tags          = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_topic" "rocketmq_topic" {
  instance_id = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  topic       = "test_topic"
  topic_type  = "NORMAL"
  queue_num   = 4
  remark      = "test for terraform"
}

`

const testAccTrocketRocketmqTopicUpdate = `

resource "tencentcloud_trocket_rocketmq_instance" "rocketmq_instance" {
  instance_type = "EXPERIMENT"
  name          = "test_topic"
  sku_code      = "experiment_500"
  remark        = "test"
  vpc_id        = "vpc-3a9fo1k9"
  subnet_id     = "subnet-8nby1yxg"
  tags          = {
    tag_key   = "rocketmq"
    tag_value = "5.x"
  }
}

resource "tencentcloud_trocket_rocketmq_topic" "rocketmq_topic" {
  instance_id = tencentcloud_trocket_rocketmq_instance.rocketmq_instance.id
  topic       = "test_topic"
  topic_type  = "NORMAL"
  queue_num   = 5
  remark      = "test terraform"
}

`
