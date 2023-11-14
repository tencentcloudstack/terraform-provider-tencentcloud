package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorAlarmPolicyDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmPolicyDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_alarm_policy.alarm_policy")),
			},
		},
	})
}

const testAccMonitorAlarmPolicyDataSource = `

data "tencentcloud_monitor_alarm_policy" "alarm_policy" {
  module = ""
  policy_name = ""
  monitor_types = 
  namespaces = 
  dimensions = ""
  receiver_uids = 
  receiver_groups = 
  policy_type = 
  field = ""
  order = ""
  project_ids = 
  notice_ids = 
  rule_types = 
  enable = 
  not_binding_notice_rule = 
  instance_group_id = 
  need_correspondence = 
  trigger_tasks {
		type = ""
		task_config = ""

  }
  one_click_policy_type = 
  not_bind_all = 
  not_instance_group = 
  prom_ins_id = ""
  receiver_on_call_form_i_ds = 
  }

`
