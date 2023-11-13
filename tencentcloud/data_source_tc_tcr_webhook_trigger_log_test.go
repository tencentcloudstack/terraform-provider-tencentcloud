package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrWebhookTriggerLogDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrWebhookTriggerLogDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcr_webhook_trigger_log.webhook_trigger_log")),
			},
		},
	})
}

const testAccTcrWebhookTriggerLogDataSource = `

data "tencentcloud_tcr_webhook_trigger_log" "webhook_trigger_log" {
  registry_id = "tcr-xxx"
  namespace = "nginx"
    tags = {
    "createdBy" = "terraform"
  }
}

`
