package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqSubscriptionAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqSubscriptionAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_subscription_attachment.subscription_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_subscription_attachment.subscription_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqSubscriptionAttachment = `

resource "tencentcloud_tdmq_subscription_attachment" "subscription_attachment" {
  environment_id = ""
  topic_name = ""
  subscription_name = ""
  is_idempotent = 
  remark = ""
  cluster_id = ""
  auto_create_policy_topic = 
  post_fix_pattern = ""
}

`
