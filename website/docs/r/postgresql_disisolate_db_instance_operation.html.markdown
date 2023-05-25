---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_disisolate_db_instance_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_disisolate_db_instance_operation"
description: |-
  Provides a resource to create a postgresql disisolate_db_instance_operation
---

# tencentcloud_postgresql_disisolate_db_instance_operation

Provides a resource to create a postgresql disisolate_db_instance_operation

## Example Usage

```hcl
resource "tencentcloud_postgresql_disisolate_db_instance_operation" "disisolate_db_instance_operation" {
  db_instance_id_set = [local.pgsql_id]
  period             = 1
  auto_voucher       = false
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id_set` - (Required, Set: [`String`], ForceNew) List of resource IDs. Note that currently you cannot remove multiple instances from isolation at the same time. Only one instance ID can be passed in here.
* `auto_voucher` - (Optional, Bool, ForceNew) Whether to use vouchers. Valid values:true (yes), false (no). Default value:false.
* `period` - (Optional, Int, ForceNew) The valid period (in months) of the monthly-subscribed instance when removing it from isolation.
* `voucher_ids` - (Optional, Set: [`String`], ForceNew) Voucher ID list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



