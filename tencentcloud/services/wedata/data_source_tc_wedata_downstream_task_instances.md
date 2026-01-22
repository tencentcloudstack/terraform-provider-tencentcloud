Use this data source to query detailed information of wedata downstream task instances

Example Usage

```hcl
data "tencentcloud_wedata_downstream_task_instances" "wedata_down_task_instances" {
  project_id   = "1859317240494305280"
  instance_key = "20250731151633120_2025-10-13 17:00:00"
}
```