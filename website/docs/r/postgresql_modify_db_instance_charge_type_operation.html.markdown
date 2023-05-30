---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_modify_db_instance_charge_type_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_modify_db_instance_charge_type_operation"
description: |-
  Provides a resource to create a postgresql modify_db_instance_charge_type_operation
---

# tencentcloud_postgresql_modify_db_instance_charge_type_operation

Provides a resource to create a postgresql modify_db_instance_charge_type_operation

## Example Usage

```hcl
resource "tencentcloud_postgresql_modify_db_instance_charge_type_operation" "modify_db_instance_charge_type_operation" {
  db_instance_id       = "postgres-6r233v55"
  instance_charge_type = "PREPAID"
  period               = 1
  auto_renew_flag      = 0
  auto_voucher         = 0
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) dbInstance ID.
* `instance_charge_type` - (Required, String, ForceNew) Instance billing mode. Valid values:PREPAID (monthly subscription), POSTPAID_BY_HOUR (pay-as-you-go).
* `period` - (Required, Int, ForceNew) Valid period in months of purchased instances. Valid values:1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. This parameter is set to 1 when the pay-as-you-go billing mode is used.
* `auto_renew_flag` - (Optional, Int, ForceNew) Renewal flag. Valid values:0 (manual renewal), 1 (auto-renewal). Default value:0.
* `auto_voucher` - (Optional, Int, ForceNew) Whether to automatically use vouchers.Valid values:1(yes),0(no).Default value:0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



