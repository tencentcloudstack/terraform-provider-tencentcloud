Use this data source to query detailed information of wedata quality rule group exec results

Example Usage

```hcl
data "tencentcloud_wedata_quality_rule_group_exec_results" "wedata_quality_rule_group_exec_results" {
  project_id = "1840731342293087232"
  filters {
    name   = "Status"
    values = ["3"]
  }
  order_fields {
    name      = "UpdateTime"
    direction = "DESC"
  }
}
```
