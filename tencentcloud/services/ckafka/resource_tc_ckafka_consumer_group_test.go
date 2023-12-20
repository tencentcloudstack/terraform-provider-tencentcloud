package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaConsumerGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
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
	instance_id = "ckafka-vv7wpvae"
	group_name = "tmp-group-name"
	topic_name_list = ["keep-topic"]
}
`
