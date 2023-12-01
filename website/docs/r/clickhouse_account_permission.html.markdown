---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_account_permission"
sidebar_current: "docs-tencentcloud-resource-clickhouse_account_permission"
description: |-
  Provides a resource to create a clickhouse account_permission
---

# tencentcloud_clickhouse_account_permission

Provides a resource to create a clickhouse account_permission

## Example Usage

```hcl
resource "tencentcloud_clickhouse_account_permission" "account_permission_all_database" {
  instance_id       = "cdwch-xxxxxx"
  cluster           = "default_cluster"
  user_name         = "user1"
  all_database      = true
  global_privileges = ["SELECT", "ALTER"]
}

resource "tencentcloud_clickhouse_account_permission" "account_permission_not_all_database" {
  instance_id  = "cdwch-xxxxxx"
  cluster      = "default_cluster"
  user_name    = "user2"
  all_database = false
  database_privilege_list {
    database_name       = "xxxxxx"
    database_privileges = ["SELECT", "ALTER"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `all_database` - (Required, Bool) Whether all database tables.
* `cluster` - (Required, String) Cluster name.
* `instance_id` - (Required, String) Instance id.
* `user_name` - (Required, String) User name.
* `database_privilege_list` - (Optional, List) Database privilege list.
* `global_privileges` - (Optional, Set: [`String`]) Global privileges.

The `database_privilege_list` object supports the following:

* `database_name` - (Required, String) Database name.
* `database_privileges` - (Optional, Set) Database privileges. Valid valuse: SELECT, INSERT_ALL, ALTER, TRUNCATE, DROP_TABLE, CREATE_TABLE, DROP_DATABASE.
* `table_privilege_list` - (Optional, List) Table privilege list.

The `table_privilege_list` object supports the following:

* `table_name` - (Required, String) Table name.
* `table_privileges` - (Required, Set) Table privileges. Valid values: SELECT, INSERT_ALL, ALTER, TRUNCATE, DROP_TABLE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

clickhouse account_permission can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_account_permission.account_permission ${instanceId}#${cluster}#${userName}
```

