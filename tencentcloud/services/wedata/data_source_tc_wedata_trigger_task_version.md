Use this data source to query detailed information of wedata trigger task version.

Example Usage

```hcl
data "tencentcloud_wedata_trigger_task_version" "trigger_task_version" {
  project_id = "1840731342175234"
  task_id    = "20241024174712123456"
  version_id = "20241024174712123456_1"
}
```
