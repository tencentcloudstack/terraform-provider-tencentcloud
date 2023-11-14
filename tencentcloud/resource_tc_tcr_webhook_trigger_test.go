package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrWebhookTriggerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrWebhookTrigger,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_webhook_trigger.webhook_trigger", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_webhook_trigger.webhook_trigger",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrWebhookTrigger = `

resource "tencentcloud_tcr_webhook_trigger" "webhook_trigger" {
  registry_id = "tcr-xxx"
  trigger {
		name = "trigger"
		targets {
			address = "http://example.org/post"
			headers {
				key = "X-Custom-Header"
				values = 
			}
		}
		event_types = 
		condition = ".*"
		enabled = true
		id = 20
		description = "this is trigger description"
		namespace_id = 10

  }
  namespace = "trigger"
  tags = {
    "createdBy" = "terraform"
  }
}

`
