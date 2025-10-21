---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_isolate_db_instance_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_isolate_db_instance_operation"
description: |-
  Provides a resource to create a postgresql isolate_db_instance_operation
---

# tencentcloud_postgresql_isolate_db_instance_operation

Provides a resource to create a postgresql isolate_db_instance_operation

## Example Usage

```hcl
resource "tencentcloud_postgresql_isolate_db_instance_operation" "isolate_db_instance_operation" {
  db_instance_id_set = [local.pgsql_id]
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id_set` - (Required, Set: [`String`], ForceNew) List of resource IDs. Note that currently you cannot isolate multiple instances at the same time. Only one instance ID can be passed in here.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



