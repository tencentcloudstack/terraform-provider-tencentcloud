---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_account_privileges"
sidebar_current: "docs-tencentcloud-resource-cynosdb_account_privileges"
description: |-
  Provides a resource to create a cynosdb account_privileges
---

# tencentcloud_cynosdb_account_privileges

Provides a resource to create a cynosdb account_privileges

## Example Usage

```hcl
resource "tencentcloud_cynosdb_account_privileges" "account_privileges" {
  cluster_id   = "cynosdbmysql-bws8h88b"
  account_name = "test"
  host         = "%"
  global_privileges = [
    "CREATE",
    "DROP",
    "ALTER",
    "CREATE TEMPORARY TABLES",
    "CREATE VIEW"
  ]
  database_privileges {
    db = "users"
    privileges = [
      "DROP",
      "REFERENCES",
      "INDEX",
      "CREATE VIEW",
      "INSERT",
      "EVENT"
    ]
  }
  table_privileges {
    db         = "users"
    table_name = "tb_user_name"
    privileges = [
      "ALTER",
      "REFERENCES",
      "SHOW VIEW"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, String, ForceNew) Account.
* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `global_privileges` - (Required, Set: [`String`]) Array of global permissions.
* `database_privileges` - (Optional, List) Array of database permissions.
* `host` - (Optional, String, ForceNew) Host, default `%`.
* `table_privileges` - (Optional, List) array of table permissions.

The `database_privileges` object supports the following:

* `db` - (Required, String) Database.
* `privileges` - (Required, Set) Database privileges.

The `table_privileges` object supports the following:

* `db` - (Required, String) Database name.
* `privileges` - (Required, Set) Table privileges.
* `table_name` - (Required, String) Table name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cynosdb account_privileges can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_account_privileges.account_privileges account_privileges_id
```

