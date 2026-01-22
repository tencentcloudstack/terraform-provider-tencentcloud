Use this data source to query detailed information of wedata trigger task versions.

Example Usage

```hcl
data "tencentcloud_wedata_trigger_task_versions" "trigger_task_versions" {
  project_id         = "1840731342175234"
  task_id           = "20241024174712123456"
  task_version_type = "SAVE"
}
```
