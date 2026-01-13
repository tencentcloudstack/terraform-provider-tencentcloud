Use this data source to query detailed information of wedata ops trigger workflow.

Example Usage

```hcl
data "tencentcloud_wedata_ops_trigger_workflow" "ops_trigger_workflow" {
  project_id  = "3108707295180644352"
  workflow_id = "b41e8d13-905a-4006-9d05-1fe180338f59"
}
```