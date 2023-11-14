package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWedataBaseline_baselineResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataBaseline_baseline,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_baseline_baseline.baseline_baseline", "id")),
			},
			{
				ResourceName:      "tencentcloud_wedata_baseline_baseline.baseline_baseline",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataBaseline_baseline = `

resource "tencentcloud_wedata_baseline_baseline" "baseline_baseline" {
  project_id = ""
  baseline_name = ""
  baseline_type = ""
  create_uin = ""
  create_name = ""
  in_charge_uin = ""
  in_charge_name = ""
  promise_tasks {
		project_id = ""
		task_name = ""
		task_id = ""
		task_cycle = ""
		workflow_name = ""
		workflow_id = ""
		task_in_charge_name = ""
		task_in_charge_uin = ""

  }
  promise_time = ""
  warning_margin = 
  is_new_alarm = 
  alarm_rule_dto {
		alarm_rule_id = ""
		alarm_level_type = ""

  }
  baseline_create_alarm_rule_request {
		project_id = ""
		creator_id = ""
		creator = ""
		rule_name = ""
		monitor_type = 
		monitor_object_ids = 
		alarm_types = 
		alarm_level = 
		alarm_ways = 
		alarm_recipient_type = 
		alarm_recipients = 
		alarm_recipient_ids = 
		ext_info = ""

  }
}

`
