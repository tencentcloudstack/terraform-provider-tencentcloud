---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_parameter_templates"
sidebar_current: "docs-tencentcloud-datasource-postgresql_parameter_templates"
description: |-
  Use this data source to query detailed information of postgresql parameter_templates
---

# tencentcloud_postgresql_parameter_templates

Use this data source to query detailed information of postgresql parameter_templates

## Example Usage

```hcl
data "tencentcloud_postgresql_parameter_templates" "parameter_templates" {
  filters {
    name   = "TemplateName"
    values = ["temp_name"]
  }
  filters {
    name   = "DBEngine"
    values = ["postgresql"]
  }
  order_by      = "CreateTime"
  order_by_type = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. Valid values:TemplateName, TemplateId, DBMajorVersion, DBEngine.
* `order_by_type` - (Optional, String) Sorting order. Valid values:asc (ascending order),desc (descending order).
* `order_by` - (Optional, String) Sorting metric. Valid values:CreateTime, TemplateName, DBMajorVersion.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) Filter name.
* `values` - (Optional, Set) One or more filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - list of parameter templates.
  * `db_engine` - the database engine for which the parameter template applies.
  * `db_major_version` - the database version to which the parameter template applies.
  * `template_description` - parameter template description.
  * `template_id` - parameter template ID.
  * `template_name` - parameter template name.


