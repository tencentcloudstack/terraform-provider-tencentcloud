package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMonitorGrafanaNotificationChannel_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaNotificationChannel,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrafanaInstanceExists("tencentcloud_monitor_grafana_notification_channel.grafanaNotificationChannel"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_notification_channel.grafanaNotificationChannel", "channel_name", "create-channel-test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_notification_channel.grafanaNotificationChannel", "org_id", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_notification_channel.grafanaNotificationChannel", "receivers.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_notification_channel.grafanaNotificationChannel", "receivers.0", "Consumer-6vkna7pevq"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_notification_channel.grafanaNotificationChannel",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaNotificationChannel = `

resource "tencentcloud_monitor_grafana_notification_channel" "grafanaNotificationChannel" {
  instance_id   = "grafana-50nj6v00"
  channel_name  = "create-channel-test"
  org_id        = 1
  receivers     = ["Consumer-6vkna7pevq"]
  extra_org_ids = []
}

`
