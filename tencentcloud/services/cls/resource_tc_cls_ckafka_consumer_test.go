package cls_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixClsCkafkaConsumerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
  compression  = 1
  need_content = true
  topic_id     = "7e34a3a7-635e-4da8-9005-88106c1fde69"
  dsl_filter   = "* | where level = 'error'"

  advanced_config {
    partition_hash_status = true
    partition_fields      = [
      "__PKGID__",
    ]
  }

  ckafka {
    instance_id   = "ckafka-qzoeaqx8"
    instance_name = "ckafka-instance"
    topic_id      = "topic-c6tm4kpm"
    topic_name    = "name"
    vip           = "172.16.112.23"
    vport         = "9092"
  }

  content {
    enable_tag         = true
    meta_fields        = [
      "__FILENAME__",
      "__HOSTNAME__",
      "__PKGID__",
      "__SOURCE__",
      "__TIMESTAMP__",
    ]
    tag_json_not_tiled = true
    timestamp_accuracy = 2
    json_type          = 1
    auto_convert_number = true
  }
}

`
