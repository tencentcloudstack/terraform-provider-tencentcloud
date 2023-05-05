package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTcrDescribeWebhookTriggerLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrDescribeWebhookTriggerLogsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcr_webhook_trigger_logs.my_logs"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "registry_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "namespace"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "trigger_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.trigger_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.event_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.notify_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.detail"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.update_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.status"),
				),
			},
		},
	})
}

const testAccTcrDescribeWebhookTriggerLogsDataSource = defaultTCRInstanceData + `

data "tencentcloud_tcr_webhook_trigger_logs" "my_logs" {
  registry_id = local.tcr_id
  namespace = var.tcr_namespace
  trigger_id = 2
    tags = {
    "createdBy" = "terraform"
  }
}

`
