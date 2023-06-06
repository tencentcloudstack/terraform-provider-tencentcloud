---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_rollback"
sidebar_current: "docs-tencentcloud-resource-mysql_rollback"
description: |-
  Provides a resource to create a mysql rollback
---

# tencentcloud_mysql_rollback

Provides a resource to create a mysql rollback

## Example Usage

```hcl
resource "tencentcloud_mysql_rollback" "rollback" {
  instance_id   = "cdb-fitq5t9h"
  strategy      = "full"
  rollback_time = "2023-05-31 23:13:35"
  databases {
    database_name     = "tf_ci_test_bak"
    new_database_name = "tf_ci_test_bak_5"
  }
  tables {
    database = "tf_ci_test_bak"
    table {
      table_name     = "test"
      new_table_name = "test_bak"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Cloud database instance ID.
* `rollback_time` - (Required, String, ForceNew) Database rollback time, the time format is: yyyy-mm-dd hh:mm:ss.
* `strategy` - (Required, String, ForceNew) Rollback strategy. Available values are: table, db, full; the default value is full. table- Extremely fast rollback mode, only import the backup and binlog of the selected table level, if there is a cross-table operation, and the associated table is not selected at the same time, the rollback will fail. In this mode, the parameter Databases must be empty; db- Quick mode, only import the backup and binlog of the selected library level, if there is a cross-database operation, and the associated library is not selected at the same time, the rollback will fail; full- normal rollback mode, the backup and binlog of the entire instance will be imported, at a slower rate.
* `databases` - (Optional, List, ForceNew) The database information to be archived, indicating that the entire database is archived.
* `tables` - (Optional, List, ForceNew) The database table information to be rolled back, indicating that the file is rolled back by table.

The `databases` object supports the following:

* `database_name` - (Required, String) The original database name before rollback.
* `new_database_name` - (Required, String) The new database name after rollback.

The `table` object supports the following:

* `new_table_name` - (Required, String) New database table name after rollback.
* `table_name` - (Required, String) The original database table name before rollback.

The `tables` object supports the following:

* `database` - (Required, String) Database name.
* `table` - (Required, List) Database table details.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



