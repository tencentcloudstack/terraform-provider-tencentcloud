package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataQualityRuleGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataQualityRuleGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_quality_rule_group.group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "project_id", "3016337760439783424"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.rule_group_name", "tf_test1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.database_name", "default"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.table_name", "big_table_500"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.monitor_type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.exec_engine_type", "HIVE"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.exec_queue", "default"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.executor_group_name", "重庆调度资源组-2a8lsema"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.description", "tf测试"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.tasks.0.task_name", "hannah_test111"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.tasks.0.workflow_id", "DATA_INTEGRATION_2025-11-01_1"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_quality_rule_group.group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataQualityRuleGroupUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_quality_rule_group.group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.rule_group_name", "tf_test1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.tasks.0.task_name", "sh_260108_205112"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.tasks.0.workflow_id", "6709198e-9bec-49ae-91ea-3a13c7160f90"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_quality_rule_group.group", "rule_group_exec_strategy_bo_list.0.tasks.0.in_charge_name_list.0", "delphichen"),
				),
			},
		},
	})
}

const testAccWedataQualityRuleGroup = `

resource "tencentcloud_wedata_quality_rule_group" "group" {
  project_id = 3016337760439783424
  rule_group_exec_strategy_bo_list {
    catalog_name        = null
    cycle_step          = 0
    cycle_type          = null
    database_name       = "default"
    datasource_id       = jsonencode(65253)
    delay_time          = 0
    description         = "tf测试"
    dlc_group_name      = null
    end_time            = null
    engine_param        = null
    exec_engine_type    = "HIVE"
    exec_plan           = null
    exec_queue          = "default"
    executor_group_id   = jsonencode(20250807142245848024)
    executor_group_name = "重庆调度资源组-2a8lsema"
    monitor_type        = 2
    rule_group_name     = "tf_test1"
    rule_id             = 0
    rule_name           = null
    schedule_time_zone  = null
    schema_name         = null
    start_time          = null
    table_name          = "big_table_500"
    task_action         = null
    trigger_types       = ["CYCLE", "MAKE_UP"]
    tasks {
      cycle_type          = 0
      in_charge_id_list   = []
      in_charge_name_list = ["hannahlliao"]
      schedule_time_zone  = null
      task_id             = jsonencode(20251118145318149)
      task_name           = "hannah_test111"
      task_type           = jsonencode(2)
      workflow_id         = "DATA_INTEGRATION_2025-11-01_1"
    }
  }
}

`

const testAccWedataQualityRuleGroupUp = `
resource "tencentcloud_wedata_quality_rule_group" "group" {
  project_id = jsonencode(3016337760439783424)
  rule_group_exec_strategy_bo_list {
    catalog_name        = null
    cycle_step          = 0
    cycle_type          = null
    database_name       = "default"
    datasource_id       = jsonencode(65253)
    delay_time          = 0
    description         = "tf测试"
    dlc_group_name      = null
    end_time            = null
    engine_param        = null
    exec_engine_type    = "HIVE"
    exec_plan           = null
    exec_queue          = "default"
    executor_group_id   = jsonencode(20250807142245848024)
    executor_group_name = "重庆调度资源组-2a8lsema"
    monitor_type        = 2
    rule_group_name     = "tf_test1"
    rule_id             = 0
    rule_name           = null
    schedule_time_zone  = null
    schema_name         = null
    start_time          = null
    table_name          = "big_table_500"
    task_action         = null
    trigger_types       = ["CYCLE", "MAKE_UP"]
    tasks {
      cycle_type          = 0
      in_charge_id_list   = ["100043191163"]
      in_charge_name_list = ["delphichen"]
      schedule_time_zone  = null
      task_id             = jsonencode(20260108210827172)
      task_name           = "sh_260108_205112"
      task_type           = jsonencode(1)
      workflow_id         = "6709198e-9bec-49ae-91ea-3a13c7160f90"
    }
  }
}


`
