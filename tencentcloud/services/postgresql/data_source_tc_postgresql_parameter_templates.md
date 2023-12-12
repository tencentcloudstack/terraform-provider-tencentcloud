Use this data source to query detailed information of postgresql parameter_templates

Example Usage

```hcl
data "tencentcloud_postgresql_parameter_templates" "parameter_templates" {
  filters {
	name = "TemplateName"
	values = ["temp_name"]
  }
  filters {
	name = "DBEngine"
	values = ["postgresql"]
  }
  order_by = "CreateTime"
  order_by_type = "desc"
}
```