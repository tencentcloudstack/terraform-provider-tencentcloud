---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_renew_db_instance_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_renew_db_instance_operation"
description: |-
  Provides a resource to create a postgresql renew_db_instance_operation
---

# tencentcloud_postgresql_renew_db_instance_operation

Provides a resource to create a postgresql renew_db_instance_operation

## Example Usage

```hcl
resource "tencentcloud_postgresql_renew_db_instance_operation" "renew_db_instance_operation" {
  db_instance_ids = "postgres-6fego161"
  period          = 12
  auto_voucher    = 0
  voucher_ids     = & lt ; nil & gt ;
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_ids` - (Required, String, ForceNew) Instance ID in the format of postgres-6fego161.
* `period` - (Required, Int, ForceNew) Renewal duration in months.
* `auto_voucher` - (Optional, Int, ForceNew) Whether to automatically use vouchers. 1:yes, 0:no. Default value:0.
* `voucher_ids` - (Optional, Set: [`String`], ForceNew) Voucher ID list (only one voucher can be specified currently).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgresql renew_db_instance_operation can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_renew_db_instance_operation.renew_db_instance_operation renew_db_instance_operation_id
```

