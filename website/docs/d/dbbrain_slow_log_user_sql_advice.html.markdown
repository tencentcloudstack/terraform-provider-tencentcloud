---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_slow_log_user_sql_advice"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_slow_log_user_sql_advice"
description: |-
  Use this data source to query detailed information of dbbrain slow_log_user_sql_advice
---

# tencentcloud_dbbrain_slow_log_user_sql_advice

Use this data source to query detailed information of dbbrain slow_log_user_sql_advice

## Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_user_sql_advice" "test" {
  instance_id = "%s"
  sql_text    = "%s"
  product     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `sql_text` - (Required, String) SQL statements.
* `product` - (Optional, String) Service product type, supported values: `mysql` - cloud database MySQL; `cynosdb` - cloud database TDSQL-C for MySQL; `dbbrain-mysql` - self-built MySQL, the default is `mysql`.
* `result_output_file` - (Optional, String) Used to save results.
* `schema` - (Optional, String) library name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `advices` - SQL optimization suggestion, which can be parsed into a JSON array, and the output is empty when no optimization is required.
* `comments` - SQL optimization suggestion remarks, which can be parsed into a String array, and the output is empty when optimization is not required.
* `cost` - The cost saving details after SQL optimization can be parsed as JSON, and the output is empty when no optimization is required.
* `request_id` - Unique request ID, returned for every request. The RequestId of the request needs to be provided when locating the problem.
* `sql_plan` - The SQL execution plan can be parsed into JSON, and the output is empty when no optimization is required.
* `tables` - The DDL information of related tables can be parsed into a JSON array.


