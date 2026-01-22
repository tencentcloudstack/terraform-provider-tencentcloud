Use this data source to query detailed information of WeData quality rule templates.

Example Usage

```hcl
data "tencentcloud_wedata_quality_rule_templates" "example" {
  project_id = "your_project_id"
  
  order_fields {
    name      = "CitationCount"
    direction = "DESC"
  }
  
  filters {
    name   = "Type"
    values = ["1", "2"]
  }
}
```
