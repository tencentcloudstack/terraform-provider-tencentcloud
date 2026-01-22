Use this data source to query detailed information of wedata upstream trigger tasks.

Example Usage

```hcl
data "tencentcloud_wedata_upstream_trigger_tasks" "upstream_trigger_tasks" {
  project_id = "3108707295180644352"
  task_id    = "20241024174712123456"
}
```