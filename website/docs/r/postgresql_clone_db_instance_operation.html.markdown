---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_clone_db_instance_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_clone_db_instance_operation"
description: |-
  Provides a resource to create a postgresql clone_db_instance_operation
---

# tencentcloud_postgresql_clone_db_instance_operation

Provides a resource to create a postgresql clone_db_instance_operation

## Example Usage

```hcl
resource "tencentcloud_postgresql_clone_db_instance_operation" "clone_db_instance_operation" {
  db_instance_id       = local.pgsql_id
  spec_code            = data.tencentcloud_postgresql_specinfos.foo.list.0.id
  storage              = data.tencentcloud_postgresql_specinfos.foo.list.0.storage_min
  period               = 1
  auto_renew_flag      = 0
  vpc_id               = local.vpc_id
  subnet_id            = local.subnet_id
  name                 = "tf_test_pg_ins_clone"
  instance_charge_type = "POSTPAID_BY_HOUR"
  security_group_ids   = [local.sg_id]
  project_id           = 0
  db_node_set {
    role = "Primary"
    zone = local.pg_az
  }
  db_node_set {
    role = "Standby"
    zone = local.pg_az
  }
  tag_list {
    tag_key   = "issued_by"
    tag_value = "terraform_test"
  }

  auto_voucher         = 0
  recovery_target_time = "%s"
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew_flag` - (Required, Int, ForceNew) Renewal flag. Valid values: `0` (manual renewal), `1` (auto-renewal). Default value: `0`.
* `db_instance_id` - (Required, String, ForceNew) ID of the original instance to be cloned.
* `db_node_set` - (Required, List, ForceNew) This parameter is required if you purchase a multi-AZ deployed instance.
* `period` - (Required, Int, ForceNew) Valid period in months of the purchased instance. Valid values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`. This parameter is set to `1` when the pay-as-you-go billing mode is used.
* `spec_code` - (Required, String, ForceNew) Purchasable specification ID, which can be obtained through the `SpecCode` field in the returned value of the `DescribeProductConfig` API.
* `storage` - (Required, Int, ForceNew) Instance storage capacity in GB.
* `subnet_id` - (Required, String, ForceNew) ID of a subnet in the VPC specified by `VpcId`.
* `vpc_id` - (Required, String, ForceNew) VPC ID.
* `activity_id` - (Optional, Int, ForceNew) Campaign ID.
* `auto_voucher` - (Optional, Int, ForceNew) Whether to automatically use vouchers. Valid values: `1` (yes), `0` (no). Default value: `0`.
* `backup_set_id` - (Optional, String, ForceNew) Basic backup set ID.
* `instance_charge_type` - (Optional, String, ForceNew) Instance billing mode. Valid values: `PREPAID` (monthly subscription), `POSTPAID_BY_HOUR` (pay-as-you-go).
* `name` - (Optional, String, ForceNew) Name of the purchased instance.
* `project_id` - (Optional, Int, ForceNew) Project ID.
* `recovery_target_time` - (Optional, String, ForceNew) Restoration point in time.
* `security_group_ids` - (Optional, Set: [`String`], ForceNew) Security group ID.
* `tag_list` - (Optional, List, ForceNew) The information of tags to be bound with the purchased instance. This parameter is left empty by default.
* `voucher_ids` - (Optional, String, ForceNew) Voucher ID list.

The `db_node_set` object supports the following:

* `role` - (Required, String) Node type. Valid values:`Primary`;`Standby`.
* `zone` - (Required, String) AZ where the node resides, such as ap-guangzhou-1.

The `tag_list` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



