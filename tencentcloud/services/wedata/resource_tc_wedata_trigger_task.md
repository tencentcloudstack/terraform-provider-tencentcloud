Provides a resource to create a wedata trigger task

Example Usage

```hcl
resource "tencentcloud_wedata_trigger_task" "trigger_task" {
  project_id = jsonencode(3108707295180644352)
  trigger_task_base_attribute {
    owner_uin        = jsonencode(100044349576)
    task_folder_path = "/"
    task_name        = "tf-test-task"
    task_type_id     = jsonencode(35)
    workflow_id      = tencentcloud_wedata_trigger_workflow.workflow.id
  }
  trigger_task_configuration {
    broker_ip           = "any"
    code_content        = base64encode("echo Hello, World")
    resource_group      = jsonencode(20241107171437783498)
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
      param_value = "eyJzcGVjTGxxxxxxxfQ=="
    }
    task_ext_configuration_list {
      param_key   = "waitExecutionTotalTTL"
      param_value = jsonencode(-1)
    }
    task_ext_configuration_list {
      param_key   = "waitExecutionTotalTTLStrategy"
      param_value = "fail"
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
```

Import

wedata trigger_task can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_trigger_task.trigger_task project_id#task_id
```