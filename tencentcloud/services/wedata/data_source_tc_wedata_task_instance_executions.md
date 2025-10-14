Use this data source to query detailed information of wedata task instance executions

Example Usage

```hcl
data "tencentcloud_wedata_task_instance_executions" "wedata_task_instance_executions" {
  project_id   = "1859317240494305280"
  instance_key = "20250731151633120_2025-10-13 17:00:00"
}
```