Use this data source to query detailed information of wedata rule templates

Example Usage

```hcl
data "tencentcloud_wedata_rule_templates" "rule_templates" {
  type                = 2
  source_object_type  = 2
  project_id          = "1840731346428280832"
  source_engine_types = [2, 4, 16]
}
```