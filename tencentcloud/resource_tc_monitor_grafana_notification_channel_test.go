package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMonitorGrafanaNotificationChannel_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGrafanaNotificationChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaNotificationChannel,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrafanaNotificationChannelExists("tencentcloud_monitor_grafana_notification_channel.grafanaNotificationChannel"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_notification_channel.grafanaNotificationChannel", "channel_name", "create-channel-test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_notification_channel.grafanaNotificationChannel", "org_id", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_notification_channel.grafanaNotificationChannel", "receivers.#", "1"),
				),
			},
		},
	})
}

func testAccCheckGrafanaNotificationChannelDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_grafana_notification_channel" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		channelId := idSplit[0]
		instanceId := idSplit[1]

		notificationChannel, err := service.DescribeMonitorGrafanaNotificationChannel(ctx, channelId, instanceId)
		if err != nil {
			return err
		}

		if notificationChannel != nil {
			return fmt.Errorf("GrafanaNotificationChannel %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckGrafanaNotificationChannelExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		channelId := idSplit[0]
		instanceId := idSplit[1]

		service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		notificationChannel, err := service.DescribeMonitorGrafanaNotificationChannel(ctx, channelId, instanceId)
		if err != nil {
			return err
		}

		if notificationChannel == nil {
			return fmt.Errorf("GrafanaNotificationChannel %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testMonitorGrafanaNotificationChannelVar = `
variable "instance_id" {
  default = "` + defaultGrafanaInstanceId + `"
}
variable "receivers" {
  default = "` + defaultGrafanaReceiver + `"
}
`

const testAccMonitorGrafanaNotificationChannel = testMonitorGrafanaNotificationChannelVar + `

resource "tencentcloud_monitor_grafana_notification_channel" "grafanaNotificationChannel" {
  instance_id   = var.instance_id
  channel_name  = "create-channel-test"
  org_id        = 1
  receivers     = [var.receivers]
  extra_org_ids = []
}

`
