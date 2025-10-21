---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_db_charsets"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_db_charsets"
description: |-
  Use this data source to query detailed information of sqlserver datasource_d_b_charsets
---

# tencentcloud_sqlserver_db_charsets

Use this data source to query detailed information of sqlserver datasource_d_b_charsets

## Example Usage

```hcl
data "tencentcloud_sqlserver_db_charsets" "example" {
  instance_id = "mssql-qelbzgwf"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID in the format of mssql-j8kv137v.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `database_charsets` - Database character set list.


