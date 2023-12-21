package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudMonitorAlarmConditionsTemplateDataSource_basic -v
func TestAccTencentCloudMonitorAlarmConditionsTemplateDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorAlarmConditionsTemplateDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.conditions.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.conditions.0.alarm_notify_period"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.conditions.0.alarm_notify_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.conditions.0.calc_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.conditions.0.calc_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.conditions.0.continue_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.conditions.0.metric_display_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.conditions.0.metric_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.conditions.0.period"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.conditions.0.rule_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.conditions.0.unit"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.insert_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.last_edit_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.update_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_alarm_conditions_template.alarm_conditions_template", "template_group_list.0.view_name"),
				),
			},
		},
	})
}

const testAccMonitorAlarmConditionsTemplateDataSource = `

data "tencentcloud_monitor_alarm_conditions_template" "alarm_conditions_template" {
  module             = "monitor"
  view_name          = "cvm_device"
  group_name         = "keep-template"
  group_id           = "7803070"
  update_time_order  = "desc=descending"
  policy_count_order = "asc=ascending"
}

`
