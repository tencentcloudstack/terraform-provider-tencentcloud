Provides a resource to create a wedata wedata_task

Example Usage

```hcl
resource "tencentcloud_wedata_workflow_folder" "wedata_workflow_folder" {
  project_id         = "2905622749543821312"
  parent_folder_path = "/"
  folder_name        = "tftest"
}

resource "tencentcloud_wedata_workflow" "wedata_workflow" {
  project_id         = 2905622749543821312
  workflow_name      = "test_workflow1"
  parent_folder_path = "${tencentcloud_wedata_workflow_folder.wedata_workflow_folder.parent_folder_path}${tencentcloud_wedata_workflow_folder.wedata_workflow_folder.folder_name}"
  workflow_type      = "cycle"
}

resource "tencentcloud_wedata_task" "wedata_task" {
  project_id = 2905622749543821312
  task_base_attribute {
    task_name = "tfTask"
    task_type_id = 30
    workflow_id = tencentcloud_wedata_workflow.wedata_workflow.workflow_id
  }
  task_configuration {
    code_content = base64encode("Hello World")
    task_ext_configuration_list {
      param_key = "bucket"
      param_value = "wedata-fusion-bjjr-1257305158"
    }
    task_ext_configuration_list {
      param_key = "ftp.file.name"
      param_value = "/datastudio/project/2905622749543821312/tftest/test_workflow1/tfTask.py"
    }
    task_ext_configuration_list {
      param_key = "tenantId"
      param_value = "1257305158"
    }
    task_ext_configuration_list {
      param_key = "region"
      param_value = "ap-beijing-fsi"
    }
    task_ext_configuration_list {
      param_key = "extraInfo"
      param_value = "{\"fromMapping\":false}"
    }
    task_ext_configuration_list {
      param_key = "ssmDynamicSkSwitch"
      param_value = "ON"
    }
    task_ext_configuration_list {
      param_key = "calendar_open"
      param_value = "0"
    }
    task_ext_configuration_list {
      param_key = "specLabelConfItems"
      param_value = "eyJzcGVjTGFiZWxDb25mSXRlbXMiOltdfQ=="
    }
    task_ext_configuration_list {
      param_key = "waitExecutionTotalTTL"
      param_value = "-1"
    }
  }
  task_scheduler_configuration {
    cycle_type = "DAY_CYCLE"
  }
}
```

Import

wedata wedata_task can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_task.wedata_task wedata_task_id
```
