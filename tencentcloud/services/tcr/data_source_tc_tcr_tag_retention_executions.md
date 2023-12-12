Use this data source to query detailed information of tcr tag_retention_executions

Example Usage

```hcl
data "tencentcloud_tcr_tag_retention_executions" "tag_retention_executions" {
  registry_id = "tcr_ins_id"
  retention_id = 1
}
```