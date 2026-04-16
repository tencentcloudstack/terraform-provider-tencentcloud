---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_base_backup"
sidebar_current: "docs-tencentcloud-resource-postgresql_base_backup"
description: |-
  Provides a resource to create a PostgreSQL base backup
---

# tencentcloud_postgresql_base_backup

Provides a resource to create a PostgreSQL base backup

## Example Usage

### Create a PostgreSQL base backup

```hcl
resource "tencentcloud_postgresql_base_backup" "example" {
  db_instance_id = "postgres-ckwcgdf1"
  tags = {
    createdBy = "Terraform"
  }
}
```

### Customize the expire time

```hcl
resource "tencentcloud_postgresql_base_backup" "example" {
  db_instance_id  = "postgres-ckwcgdf1"
  new_expire_time = "2027-04-23 20:07:36"
  tags = {
    createdBy = "Terraform"
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


## Import

PostgreSQL base backup can be imported using the dBInstanceId#baseBackupId, e.g.

```
terraform import tencentcloud_postgresql_base_backup.example postgres-ckwcgdf1#bac3d001-5160-5077-9139-49c1310e0854
```

