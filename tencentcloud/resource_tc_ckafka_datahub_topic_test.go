package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCkafkaDatahubTopicResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaDatahubTopic,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ckafka_datahub_topic.datahub_topic", "id")),
			},
			{
				Config: testAccCkafkaDatahubTopicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_datahub_topic.datahub_topic", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_datahub_topic.datahub_topic", "retention_ms", "120000"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_datahub_topic.datahub_topic", "note", "for test 123"),
				),
			},
			{
				ResourceName:      "tencentcloud_ckafka_datahub_topic.datahub_topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCkafkaDatahubTopic = `

data "tencentcloud_user_info" "user" {}

resource "tencentcloud_ckafka_datahub_topic" "datahub_topic" {
  name = format("%s-tf", data.tencentcloud_user_info.user.app_id)
  partition_num = 20
  retention_ms = 60000
  note = "for test"
  tags = {
    "createdBy" = "terraform"
  }
}

`

const testAccCkafkaDatahubTopicUpdate = `

data "tencentcloud_user_info" "user" {}

resource "tencentcloud_ckafka_datahub_topic" "datahub_topic" {
  name = format("%s-tf", data.tencentcloud_user_info.user.app_id)
  partition_num = 20
  retention_ms = 120000
  note = "for test 123"
  tags = {
    "createdBy" = "terraform"
  }
}

`
