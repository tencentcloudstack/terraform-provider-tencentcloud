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
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_notification_channel.grafanaNotificationChannel", "id"),
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
  instance_id = ""
    channel_name = ""
  org_id = 1
  receivers = ""
  extra_org_ids = ""
}

`
