Use this data source to query detailed information of wedata downstream trigger tasks.

Example Usage

```hcl
data "tencentcloud_wedata_downstream_trigger_tasks" "downstream_trigger_tasks" {
  project_id = "3108707295180644352"
  task_id    = "20241024174712123456"
}
```