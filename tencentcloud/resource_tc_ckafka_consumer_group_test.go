package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaConsumerGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaConsumerGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ckafka_consumer_group.consumer_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_ckafka_consumer_group.consumer_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCkafkaConsumerGroup = `

resource "tencentcloud_ckafka_consumer_group" "consumer_group" {
  instance_id = "InstanceId"
  group_name = "GroupName"
  topic_name = "TopicName"
  topic_name_list = 
}

`
