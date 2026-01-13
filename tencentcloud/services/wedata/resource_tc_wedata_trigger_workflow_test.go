package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTriggerWorkflowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTriggerWorkflow,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_trigger_workflow.workflow", "id"),
					// Required fields
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "project_id", "3107469419838337024"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "workflow_name", "tf-test1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "parent_folder_path", "/默认文件夹"),
					// Optional fields
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "owner_uin", "100044349576"),
					// Workflow parameters
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "workflow_params.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "workflow_params.0.param_key", "aaa"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "workflow_params.0.param_value", "bbb"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "workflow_params.1.param_key", "bbb"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "workflow_params.1.param_value", "ccc"),
					// General task parameters
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "general_task_params.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "general_task_params.0.type", "SPARK_SQL"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "general_task_params.0.value", "a=b\nb=c\nc=d\nd=e"),
					// Trigger workflow scheduler configurations
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.trigger_mode", "TIME_TRIGGER"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.config_mode", "COMMON"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.cycle_type", "DAY_CYCLE"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.crontab_expression", "0 0 * * * ? *"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.schedule_time_zone", "UTC+8"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.scheduler_status", "ACTIVE"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.start_time", "2026-01-09 00:00:00"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.end_time", "2099-12-31 23:59:59"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.trigger_minimum_interval_second", "0"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.trigger_wait_time_second", "0"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.extra_info", ""),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.file_arrival_path", ""),
					// Computed fields
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.trigger_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_trigger_workflow.workflow",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataTriggerWorkflowUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_trigger_workflow.workflow", "id"),
					// Required fields
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "project_id", "3107469419838337024"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "workflow_name", "tf-test1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "parent_folder_path", "/默认文件夹"),
					// Optional fields
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "owner_uin", "100044349576"),
					// Workflow parameters (updated - only 1 param now)
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "workflow_params.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "workflow_params.0.param_key", "aaa"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "workflow_params.0.param_value", "bbb"),
					// General task parameters (updated value)
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "general_task_params.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "general_task_params.0.type", "SPARK_SQL"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "general_task_params.0.value", "a=b\nb=c"),
					// Trigger workflow scheduler configurations (updated values)
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.trigger_mode", "TIME_TRIGGER"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.config_mode", "COMMON"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.cycle_type", "DAY_CYCLE"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.crontab_expression", "0 0 0 * * ? *"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.schedule_time_zone", "UTC+8"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.scheduler_status", "ACTIVE"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.start_time", "2026-01-07 00:00:00"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.end_time", "2099-12-31 23:59:59"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.trigger_minimum_interval_second", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.trigger_wait_time_second", "1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.extra_info", ""),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.file_arrival_path", ""),
					// Computed fields
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_trigger_workflow.workflow", "trigger_workflow_scheduler_configurations.0.trigger_id"),
				),
			},
		},
	})
}

const testAccWedataTriggerWorkflow = `
resource "tencentcloud_wedata_trigger_workflow" "workflow" {
  bundle_id          = null
  bundle_info        = null
  owner_uin          = 100044349576
  parent_folder_path = "/默认文件夹"
  project_id         = 3107469419838337024
  workflow_desc      = null
  workflow_name      = "tf-test1"
  general_task_params {
    type  = "SPARK_SQL"
    value = "a=b\nb=c\nc=d\nd=e"
  }
  trigger_workflow_scheduler_configurations {
    config_mode                     = "COMMON"
    crontab_expression              = "0 0 * * * ? *"
    cycle_type                      = "DAY_CYCLE"
    end_time                        = "2099-12-31 23:59:59"
    extra_info                      = null
    file_arrival_path               = null
    schedule_time_zone              = "UTC+8"
    scheduler_status                = "ACTIVE"
    start_time                      = "2026-01-09 00:00:00"
    trigger_minimum_interval_second = 0
    trigger_mode                    = "TIME_TRIGGER"
    trigger_wait_time_second        = 0
  }
  workflow_params {
    param_key   = "aaa"
    param_value = "bbb"
  }
  workflow_params {
    param_key   = "bbb"
    param_value = "ccc"
  }
}
`

const testAccWedataTriggerWorkflowUp = `
resource "tencentcloud_wedata_trigger_workflow" "workflow" {
  bundle_id          = null
  bundle_info        = null
  owner_uin          = 100044349576
  parent_folder_path = "/默认文件夹"
  project_id         = 3107469419838337024
  workflow_desc      = null
  workflow_name      = "tf-test1"
  general_task_params {
    type  = "SPARK_SQL"
    value = "a=b\nb=c"
  }
  trigger_workflow_scheduler_configurations {
    config_mode                     = "COMMON"
    crontab_expression              = "0 0 0 * * ? *"
    cycle_type                      = "DAY_CYCLE"
    end_time                        = "2099-12-31 23:59:59"
    extra_info                      = null
    file_arrival_path               = null
    schedule_time_zone              = "UTC+8"
    scheduler_status                = "ACTIVE"
    start_time                      = "2026-01-07 00:00:00"
    trigger_minimum_interval_second = 1
    trigger_mode                    = "TIME_TRIGGER"
    trigger_wait_time_second        = 1
  }
  workflow_params {
    param_key   = "aaa"
    param_value = "bbb"
  }
}
`
