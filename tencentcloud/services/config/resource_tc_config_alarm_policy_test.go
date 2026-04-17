package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigAlarmPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAlarmPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_config_alarm_policy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_config_alarm_policy.example", "alarm_policy_id"),
					resource.TestCheckResourceAttr("tencentcloud_config_alarm_policy.example", "name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_config_alarm_policy.example", "status", "1"),
				),
			},
			{
				Config: testAccConfigAlarmPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_config_alarm_policy.example", "name", "tf-example-update"),
					resource.TestCheckResourceAttr("tencentcloud_config_alarm_policy.example", "status", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_config_alarm_policy.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccConfigAlarmPolicy = `
resource "tencentcloud_config_alarm_policy" "example" {
  name                   = "tf-example"
  event_scope            = [1]
  risk_level             = [1, 2]
  notice_time            = "09:30:00~23:30:00"
  notification_mechanism = "实时发送"
  status                 = 1
  notice_period          = [1, 2, 3, 4, 5]
  description            = "tf example alarm policy"
}
`

const testAccConfigAlarmPolicyUpdate = `
resource "tencentcloud_config_alarm_policy" "example" {
  name                   = "tf-example-update"
  event_scope            = [1]
  risk_level             = [1, 2, 3]
  notice_time            = "00:00:00~23:59:59"
  notification_mechanism = "实时发送"
  status                 = 2
  notice_period          = [1, 2, 3, 4, 5, 6, 7]
  description            = "updated alarm policy"
}
`
