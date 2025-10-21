---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_base_backup"
sidebar_current: "docs-tencentcloud-resource-postgresql_base_backup"
description: |-
  Provides a resource to create a postgresql base_backup
---

# tencentcloud_postgresql_base_backup

Provides a resource to create a postgresql base_backup

## Example Usage

```hcl
resource "tencentcloud_postgresql_base_backup" "base_backup" {
  db_instance_id = local.pgsql_id
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String) Instance ID.
* `new_expire_time` - (Optional, String) New expiration time.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `base_backup_id` - Base backup ID.


