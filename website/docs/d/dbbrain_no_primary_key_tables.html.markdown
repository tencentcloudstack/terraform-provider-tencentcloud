---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_no_primary_key_tables"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_no_primary_key_tables"
description: |-
  Use this data source to query detailed information of dbbrain no_primary_key_tables
---

# tencentcloud_dbbrain_no_primary_key_tables

Use this data source to query detailed information of dbbrain no_primary_key_tables

## Example Usage

```hcl
data "tencentcloud_dbbrain_no_primary_key_tables" "no_primary_key_tables" {
  instance_id = ""
  date        = ""
  product     = ""
}
```

## Argument Reference

The following arguments are supported:

* `date` - (Required, String) Query date, such as 2021-05-27, the earliest date is 30 days ago.
* `instance_id` - (Required, String) instance id.
* `product` - (Optional, String) Service product type, supported values: `mysql` - ApsaraDB for MySQL, the default is `mysql`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `no_primary_key_table_count_diff` - The difference with yesterday&amp;#39;s scan of the table without a primary key. A positive number means an increase, a negative number means a decrease, and 0 means no change.
* `no_primary_key_tables` - A list of tables without primary keys.
  * `engine` - Storage engine for database tables.
  * `table_name` - tableName.
  * `table_rows` - rows.
  * `table_schema` - library name.
  * `total_length` - Total space used (MB).
* `timestamp` - Collection timestamp (seconds).


