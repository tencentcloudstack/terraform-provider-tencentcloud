package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorGrafanaNotificationChannelResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaNotificationChannel,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_notification_channel.grafana_notification_channel", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_notification_channel.grafana_notification_channel",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorGrafanaNotificationChannel = `

resource "tencentcloud_monitor_grafana_notification_channel" "grafana_notification_channel" {
  instance_id = &lt;nil&gt;
    channel_name = &lt;nil&gt;
  org_id = 1
  receivers = &lt;nil&gt;
  extra_org_ids = &lt;nil&gt;
}

`
