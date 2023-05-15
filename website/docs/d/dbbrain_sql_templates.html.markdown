---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_sql_templates"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_sql_templates"
description: |-
  Use this data source to query detailed information of dbbrain sql_templates
---

# tencentcloud_dbbrain_sql_templates

Use this data source to query detailed information of dbbrain sql_templates

## Example Usage

```hcl
data "tencentcloud_dbbrain_sql_templates" "sql_templates" {
  instance_id = ""
  schema      = ""
  sql_text    = ""
  product     = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `schema` - (Required, String) database name.
* `sql_text` - (Required, String) SQL statements.
* `product` - (Optional, String) Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `sql_id` - SQL template ID.
* `sql_template` - SQL template content.
* `sql_type` - sql type.


