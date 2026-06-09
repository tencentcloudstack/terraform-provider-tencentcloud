package cls_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClsKafkaConsumerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsKafkaConsumer,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_kafka_consumer.example", "from_topic_id"),
				),
			},
			{
				Config: testAccClsKafkaConsumerUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_kafka_consumer.example", "from_topic_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_kafka_consumer.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsKafkaConsumer = `
resource "tencentcloud_cls_kafka_consumer" "example" {
  from_topic_id = "c9b68233-948a-4eaf-a363-d0c2ced393ae"
  compression   = 0
  consumer_content {
    enable_tag      = false
    format          = 1
    json_type       = 1
    meta_fields     = [
      "__SOURCE__",
      "__FILENAME__"
    ]
    tag_transaction = 2
  }
}
`

const testAccClsKafkaConsumerUpdate = `
resource "tencentcloud_cls_kafka_consumer" "example" {
  from_topic_id = "c9b68233-948a-4eaf-a363-d0c2ced393ae"
  compression   = 0
  consumer_content {
    enable_tag      = false
    format          = 0
    json_type       = 1
    meta_fields     = [
      "__SOURCE__",
    ]
    tag_transaction = 2
  }
}
`
