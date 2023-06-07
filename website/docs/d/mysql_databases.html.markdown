---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_databases"
sidebar_current: "docs-tencentcloud-datasource-mysql_databases"
description: |-
  Use this data source to query detailed information of mysql databases
---

# tencentcloud_mysql_databases

Use this data source to query detailed information of mysql databases

## Example Usage

```hcl
data "tencentcloud_mysql_databases" "databases" {
  instance_id     = "cdb-c1nl9rpv"
  database_regexp = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of instance.
* `database_regexp` - (Optional, String) Regular expression to match database library names.
* `limit` - (Optional, Int) The number of single requests, the default value is 20, the minimum value is 1, and the maximum value is 100.
* `offset` - (Optional, Int) Page offset.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `database_list` - Database name and character set.
  * `character_set` - character set type.
  * `database_name` - The name of database.
* `items` - Returned instance information.


