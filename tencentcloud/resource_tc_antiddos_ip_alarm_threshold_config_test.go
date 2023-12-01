package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosIpAlarmThresholdConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosIpAlarmThresholdConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_ip_alarm_threshold_config.ip_alarm_threshold_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ip_alarm_threshold_config.ip_alarm_threshold_config", "alarm_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ip_alarm_threshold_config.ip_alarm_threshold_config", "alarm_threshold", "2"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ip_alarm_threshold_config.ip_alarm_threshold_config", "instance_ip", "119.91.67.112"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ip_alarm_threshold_config.ip_alarm_threshold_config", "instance_id", "bgp-00000ry7"),
				),
			},
			{
				Config: testAccAntiddosIpAlarmThresholdConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_ip_alarm_threshold_config.ip_alarm_threshold_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_ip_alarm_threshold_config.ip_alarm_threshold_config", "alarm_threshold", "3"),
				),
			},
			{
				ResourceName:      "tencentcloud_antiddos_ip_alarm_threshold_config.ip_alarm_threshold_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosIpAlarmThresholdConfig = `
resource "tencentcloud_antiddos_ip_alarm_threshold_config" "ip_alarm_threshold_config" {
	alarm_type = 1
	alarm_threshold = 2
	instance_ip = "119.91.67.112"
	instance_id = "bgp-00000ry7"
}
`

const testAccAntiddosIpAlarmThresholdConfigUpdate = `
resource "tencentcloud_antiddos_ip_alarm_threshold_config" "ip_alarm_threshold_config" {
	alarm_type = 1
	alarm_threshold = 3
	instance_ip = "119.91.67.112"
	instance_id = "bgp-00000ry7"
}
`
