package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataTriggerTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataTriggerTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_trigger_task.trigger_task", "id"),
					// Required fields
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "project_id", "3108707295180644352"),
					// Basic attributes
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_base_attribute.0.owner_uin", "100044349576"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_base_attribute.0.task_folder_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_base_attribute.0.task_name", "tf-test-task"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_base_attribute.0.task_type_id", "35"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_base_attribute.0.workflow_id", "7b716065-d016-41aa-bc68-e18a1253df7e"),
					// Configuration
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_configuration.0.broker_ip", "any"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_configuration.0.resource_group", "20241107171437783498"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_configuration.0.code_content"),
					// Extended configuration count
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_configuration.0.task_ext_configuration_list.#", "7"),
					// Scheduler configuration
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.allow_redo_type", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.execution_ttl_minute", "-1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.max_retry_number", "4"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.retry_wait_minute", "5"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.run_priority_type", "6"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.wait_execution_total_ttl_minute", "-1"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_trigger_task.trigger_task",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataTriggerTaskUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_trigger_task.trigger_task", "id"),
					// Required fields
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "project_id", "3108707295180644352"),
					// Basic attributes (updated)
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_base_attribute.0.owner_uin", "100044349576"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_base_attribute.0.task_folder_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_base_attribute.0.task_name", "tf-test-task1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_base_attribute.0.task_type_id", "35"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_base_attribute.0.workflow_id", "7b716065-d016-41aa-bc68-e18a1253df7e"),
					// Configuration (updated)
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_configuration.0.broker_ip", "any"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_configuration.0.resource_group", "20241107171437783498"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_configuration.0.code_content"),
					// Extended configuration count (updated)
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_configuration.0.task_ext_configuration_list.#", "7"),
					// Scheduler configuration (same)
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.allow_redo_type", "ALL"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.execution_ttl_minute", "-1"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.max_retry_number", "5"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.retry_wait_minute", "5"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.run_priority_type", "6"),
					resource.TestCheckResourceAttr("tencentcloud_wedata_trigger_task.trigger_task", "trigger_task_scheduler_configuration.0.wait_execution_total_ttl_minute", "-1"),
				),
			},
		},
	})
}

const testAccWedataTriggerTask = `


resource "tencentcloud_wedata_trigger_task" "trigger_task" {
  project_id = 3108707295180644352
  trigger_task_base_attribute {
    owner_uin        = 100044349576
    task_description = null
    task_folder_path = "/"
    task_name        = "tf-test-task"
    task_type_id     = 35
    workflow_id      = "7b716065-d016-41aa-bc68-e18a1253df7e"
  }
  trigger_task_configuration {
    broker_ip           = "any"
    bundle_id           = null
    bundle_info         = null
    code_content        = base64encode("echo Hello, World")
    data_cluster        = null
    resource_group      = 20241107171437783498
    source_service_id   = null
    source_service_name = null
    source_service_type = null
    target_service_id   = null
    target_service_name = null
    target_service_type = null
    yarn_queue          = null
    task_ext_configuration_list {
      param_key   = "enableKerberosLogin"
      param_value = false
    }
    task_ext_configuration_list {
      param_key   = "executionTTLStrategy"
      param_value = "fail"
    }
    task_ext_configuration_list {
      param_key   = "python_sub_version"
      param_value = "python3"
    }
    task_ext_configuration_list {
      param_key   = "python_type"
      param_value = "python3"
    }
    task_ext_configuration_list {
      param_key   = "specLabelConfItems"
      param_value = "eyJzcGVjTGFiZWxDb25mSXRlbXMiOltdfQ=="
    }
    task_ext_configuration_list {
      param_key   = "waitExecutionTotalTTL"
      param_value = jsonencode(-1)
    }
    task_ext_configuration_list {
      param_key   = "waitExecutionTotalTTLStrategy"
      param_value = "fail"
    }
  }
  trigger_task_scheduler_configuration {
    allow_redo_type                 = "ALL"
    execution_ttl_minute            = -1
    max_retry_number                = 4
    retry_wait_minute               = 5
    run_priority_type               = 6
    wait_execution_total_ttl_minute = -1
  }
}


`

const testAccWedataTriggerTaskUp = `

resource "tencentcloud_wedata_trigger_task" "trigger_task" {
  project_id = 3108707295180644352
  trigger_task_base_attribute {
    owner_uin        = 100044349576
    task_description = null
    task_folder_path = "/"
    task_name        = "tf-test-task1"
    task_type_id     = 35
    workflow_id      = "7b716065-d016-41aa-bc68-e18a1253df7e"
  }
  trigger_task_configuration {
    broker_ip           = "any"
    bundle_id           = null
    bundle_info         = null
    code_content        = base64encode("echo Hello, World!")
    data_cluster        = null
    resource_group      = 20241107171437783498
    source_service_id   = null
    source_service_name = null
    source_service_type = null
    target_service_id   = null
    target_service_name = null
    target_service_type = null
    yarn_queue          = null
    task_ext_configuration_list {
      param_key   = "enableKerberosLogin"
      param_value = true
    }
    task_ext_configuration_list {
      param_key   = "executionTTLStrategy"
      param_value = "fail"
    }
    task_ext_configuration_list {
      param_key   = "python_sub_version"
      param_value = "python3"
    }
    task_ext_configuration_list {
      param_key   = "python_type"
      param_value = "python3"
    }
    task_ext_configuration_list {
      param_key   = "specLabelConfItems"
      param_value = "eyJzcGVjTGFiZWxDb25mSXRlbXMiOltdfQ=="
    }
    task_ext_configuration_list {
      param_key   = "waitExecutionTotalTTL"
      param_value = jsonencode(-1)
    }
    task_ext_configuration_list {
      param_key   = "waitExecutionTotalTTLStrategy"
      param_value = "fail"
    }
  }
  trigger_task_scheduler_configuration {
    allow_redo_type                 = "ALL"
    execution_ttl_minute            = -1
    max_retry_number                = 5
    retry_wait_minute               = 5
    run_priority_type               = 6
    wait_execution_total_ttl_minute = -1
  }
}


`
