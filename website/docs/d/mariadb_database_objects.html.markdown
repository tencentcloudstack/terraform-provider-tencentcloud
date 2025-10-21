---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_database_objects"
sidebar_current: "docs-tencentcloud-datasource-mariadb_database_objects"
description: |-
  Use this data source to query detailed information of mariadb database_objects
---

# tencentcloud_mariadb_database_objects

Use this data source to query detailed information of mariadb database_objects

## Example Usage

```hcl
data "tencentcloud_mariadb_database_objects" "database_objects" {
  instance_id = "tdsql-n2fw7pn3"
  db_name     = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `db_name` - (Required, String) database name.
* `instance_id` - (Required, String) instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `funcs` - func list.
  * `func` - func name.
* `procs` - proc list.
  * `proc` - proc name.
* `tables` - table list.
  * `table` - table name.
* `views` - view list.
  * `view` - view name.


