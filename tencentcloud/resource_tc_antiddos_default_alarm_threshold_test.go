package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosDefaultAlarmThresholdResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosDefaultAlarmThreshold,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold", "default_alarm_config.0.alarm_threshold", "1"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold", "default_alarm_config.0.alarm_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold", "instance_type", "bgp"),
				),
			},
			{
				Config: testAccAntiddosDefaultAlarmThresholdUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold", "id"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold", "default_alarm_config.0.alarm_threshold", "2"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold", "default_alarm_config.0.alarm_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold", "instance_type", "bgp"),
				),
			},
			{
				ResourceName:      "tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosDefaultAlarmThreshold = `
resource "tencentcloud_antiddos_default_alarm_threshold" "default_alarm_threshold" {
	default_alarm_config {
	  alarm_type = 1
	  alarm_threshold = 1
	}
	instance_type = "bgp"
  }
`

const testAccAntiddosDefaultAlarmThresholdUpdate = `
resource "tencentcloud_antiddos_default_alarm_threshold" "default_alarm_threshold" {
	default_alarm_config {
	  alarm_type = 1
	  alarm_threshold = 2
	}
	instance_type = "bgp"
  }
`
