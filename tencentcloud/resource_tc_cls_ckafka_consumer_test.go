package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClsCkafkaConsumerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsCkafkaConsumer,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cls_ckafka_consumer.ckafka_consumer", "id")),
			},
			{
				ResourceName:      "tencentcloud_cls_ckafka_consumer.ckafka_consumer",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsCkafkaConsumer = `

resource "tencentcloud_cls_ckafka_consumer" "ckafka_consumer" {
  topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  need_content = true
  content {
		enable_tag = true
		meta_fields = 
		tag_json_not_tiled = true
		timestamp_accuracy = 1

  }
  ckafka {
		vip = "1.1.1.1"
		vport = "8000"
		instance_id = "ckafka-xxxxx"
		instance_name = "test"
		topic_id = "topic-5cd3a17e-fb0b-418c-afd7-77b3653974xx"
		topic_name = "test"

  }
  compression = 0
}

`
