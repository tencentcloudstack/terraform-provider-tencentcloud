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
  db_instance_id = tencentcloud_postgresql_instance.oper_test_PREPAID.id
  period         = 1
  auto_voucher   = 0
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance ID in the format of postgres-6fego161.
* `period` - (Required, Int, ForceNew) Renewal duration in months.
* `auto_voucher` - (Optional, Int, ForceNew) Whether to automatically use vouchers. 1:yes, 0:no. Default value:0.
* `voucher_ids` - (Optional, Set: [`String`], ForceNew) Voucher ID list (only one voucher can be specified currently).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



