package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorAlarmPolicyDataSource_basic -v
func TestAccTencentCloudMonitorAlarmPolicyDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmPolicyDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_alarm_policy.alarm_policy"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.can_set_default"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.0.continue_period"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.0.is_open"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.0.metric_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.0.notice_frequency"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.0.operator"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.0.period"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.0.rule_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.0.unit"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.0.value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.condition.0.rules.0.value_max"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.conditions_temp.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.enable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.event_condition.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.event_condition.0.rules.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.event_condition.0.rules.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.event_condition.0.rules.0.metric_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.insert_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.last_edit_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.monitor_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.namespace"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.namespace_show_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.notice_ids.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.notices.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.notices.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.notices.0.is_preset"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.notices.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.notices.0.notice_language"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.notices.0.notice_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.notices.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.notices.0.updated_by"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.notices.0.user_notices.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.origin_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.policy_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.policy_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.project_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.region.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.rule_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.tag_instances.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.tag_instances.0.binding_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.tag_instances.0.key"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.tag_instances.0.region_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.tag_instances.0.service_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.tag_instances.0.tag_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.tag_instances.0.value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_policy.alarm_policy", "policies.0.update_time"),
				),
			},
		},
	})
}

const testAccMonitorAlarmPolicyDataSource = `

data "tencentcloud_monitor_alarm_policy" "alarm_policy" {
  module        = "monitor"
  policy_name   = "terraform"
  monitor_types = ["MT_QCE"]
  namespaces    = ["cvm_device"]
  project_ids   = [0]
  notice_ids    = ["notice-f2svbu3w"]
  rule_types    = ["STATIC"]
  enable        = [1]
}

`
