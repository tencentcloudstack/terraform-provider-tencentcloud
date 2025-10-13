package wedata_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"testing"
)

func TestAccTencentCloudWedataOpsAlarmRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataOpsAlarmRule,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule", "id")),
		}, {
			ResourceName:      "tencentcloud_wedata_ops_alarm_rule.wedata_ops_alarm_rule",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWedataOpsAlarmRule = `

resource "tencentcloud_wedata_ops_alarm_rule" "wedata_ops_alarm_rule" {
  alarm_groups = {
    notification_fatigue = {
      quiet_intervals = {
      }
    }
    web_hooks = {
    }
  }
  alarm_rule_detail = {
    time_out_ext_info = {
    }
    data_backfill_or_rerun_time_out_ext_info = {
    }
    project_instance_statistics_alarm_info_list = {
    }
    reconciliation_ext_info = {
    }
  }
}
`
