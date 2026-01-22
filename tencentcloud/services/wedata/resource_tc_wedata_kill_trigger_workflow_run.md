Provides a resource to kill wedata trigger workflow run

~> **NOTE:** Both "all" and "pending" require obtaining the execution_id through the query interface before passing it as a parameter..

Example Usage

```hcl

data "tencentcloud_wedata_trigger_workflow_runs" "trigger_workflow_runs" {
  project_id = "3108707295180644352"
  filters {
    name   = "WorkflowId"
    values = ["333368d7-bc8e-4b95-9a66-7a5151063eb2"]
  }
  order_fields {
    name      = "CreateTime"
    direction = "DESC"
  }
}

resource "tencentcloud_wedata_kill_trigger_workflow_run" "kill_all" {
  project_id  = "3108707295180644352"
  workflow_id = "333368d7-bc8e-4b95-9a66-7a5151063eb2"
  workflow_execution_id= data.tencentcloud_wedata_trigger_workflow_runs.trigger_workflow_runs.data[0].items[0].execution_id
}
```