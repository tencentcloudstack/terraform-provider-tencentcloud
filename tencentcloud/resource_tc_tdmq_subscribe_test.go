package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqSubscribeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqSubscribe,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_subscribe.subscribe", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_subscribe.subscribe",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqSubscribe = `

resource "tencentcloud_tdmq_subscribe" "subscribe" {
  topic_name = "topic_name"
  subscription_name = "subscription_name"
  protocol = "HTTP"
  endpoint = &lt;nil&gt;
  notify_strategy = "EXPONENTIAL_DECAY_RETRY"
  filter_tag = &lt;nil&gt;
  binding_key = &lt;nil&gt;
  notify_content_format = "JSON"
  tags = {
    "createdBy" = "terraform"
  }
}

`
