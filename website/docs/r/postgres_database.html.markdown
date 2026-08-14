---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgres_database"
sidebar_current: "docs-tencentcloud-resource-postgres_database"
description: |-
  Provides a resource to manage a TencentDB for PostgreSQL database within a DB instance.
---

# tencentcloud_postgres_database

Provides a resource to manage a TencentDB for PostgreSQL database within a DB instance.

## Example Usage

```hcl
resource "tencentcloud_postgres_database" "example" {
  db_instance_id = "postgres-6fego161"
  database_name  = "test_db"
  database_owner = "tcuser"
  encoding       = "UTF8"
  collate        = "C"
  ctype          = "C"
}
```

## Argument Reference

The following arguments are supported:

* `database_name` - (Required, String, ForceNew) Database name.
* `database_owner` - (Required, String) Database owner account.
* `db_instance_id` - (Required, String, ForceNew) DB instance ID, such as postgres-6fego161.
* `collate` - (Optional, String, ForceNew) Database collation rule.
* `ctype` - (Optional, String, ForceNew) Database character classification.
* `encoding` - (Optional, String, ForceNew) Database character encoding, such as UTF8, LATIN1, LATIN2, WIN1250, WIN1251, WIN1252, KOI8R, EUC_JP, EUC_KR. Default: UTF8.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

PostgreSQL database can be imported using the composite ID `db_instance_id#database_name`, e.g.

```
terraform import tencentcloud_postgres_database.example postgres-6fego161#test_db
```

