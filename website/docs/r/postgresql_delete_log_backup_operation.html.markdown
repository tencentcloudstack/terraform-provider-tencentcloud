---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_delete_log_backup_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_delete_log_backup_operation"
description: |-
  Provides a resource to create a postgresql delete_log_backup_operation
---

# tencentcloud_postgresql_delete_log_backup_operation

Provides a resource to create a postgresql delete_log_backup_operation

## Example Usage

```hcl
resource "tencentcloud_postgresql_delete_log_backup_operation" "delete_log_backup_operation" {
  db_instance_id = "local.pg_id"
  log_backup_id  = "local.pg_log_backup_id"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance ID.
* `log_backup_id` - (Required, String, ForceNew) Log backup ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



