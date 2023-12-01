package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorAlarmPolicySetDefaultResource_basic -v
func TestAccTencentCloudMonitorAlarmPolicySetDefaultResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmPolicySetDefault,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_alarm_policy_set_default.policy_set_default", "id"),
				),
			},
		},
	})
}

const testAccMonitorAlarmPolicySetDefault = testAccMonitorAlarmPolicy + `

resource "tencentcloud_monitor_alarm_policy_set_default" "policy_set_default" {
  module = "monitor"
  policy_id = tencentcloud_monitor_alarm_policy.policy.id
}

`
