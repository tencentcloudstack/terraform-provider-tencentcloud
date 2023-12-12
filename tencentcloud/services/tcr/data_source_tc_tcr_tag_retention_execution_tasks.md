Use this data source to query detailed information of tcr tag_retention_execution_tasks

Example Usage

```hcl
data "tencentcloud_tcr_tag_retention_execution_tasks" "tasks" {
  registry_id = "tcr_ins_id"
  retention_id = 1
  execution_id = 1
}
```