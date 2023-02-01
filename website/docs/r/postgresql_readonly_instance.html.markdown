---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_readonly_instance"
sidebar_current: "docs-tencentcloud-resource-postgresql_readonly_instance"
description: |-
  Use this resource to create postgresql readonly instance.
---

# tencentcloud_postgresql_readonly_instance

Use this resource to create postgresql readonly instance.

## Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_instance" "foo" {
  auto_renew_flag       = 0
  db_version            = "10.4"
  instance_charge_type  = "POSTPAID_BY_HOUR"
  master_db_instance_id = "postgres-j4pm65id"
  memory                = 4
  name                  = "hello"
  need_support_ipv6     = 0
  project_id            = 0
  security_groups_ids = [
    "sg-fefj5n6r",
  ]
  storage   = 250
  subnet_id = "subnet-enm92y0m"
  vpc_id    = "vpc-86v957zb"
  zone      = "ap-guangzhou-6"
}
```

## Argument Reference

The following arguments are supported:

* `db_version` - (Required, String, ForceNew) PostgreSQL kernel version, which must be the same as that of the primary instance.
* `master_db_instance_id` - (Required, String, ForceNew) ID of the primary instance to which the read-only replica belongs.
* `memory` - (Required, Int) Memory size(in GB). Allowed value must be larger than `memory` that data source `tencentcloud_postgresql_specinfos` provides.
* `name` - (Required, String) Instance name.
* `project_id` - (Required, Int) Project ID.
* `security_groups_ids` - (Required, Set: [`String`]) ID of security group.
* `storage` - (Required, Int) Instance storage capacity in GB.
* `subnet_id` - (Required, String) VPC subnet ID.
* `vpc_id` - (Required, String, ForceNew) VPC ID.
* `zone` - (Required, String, ForceNew) Availability zone ID, which can be obtained through the Zone field in the returned value of the DescribeZones API.
* `auto_renew_flag` - (Optional, Int) Auto renew flag, `1` for enabled. NOTES: Only support prepaid instance.
* `auto_voucher` - (Optional, Int) Whether to use voucher, `1` for enabled.
* `instance_charge_type` - (Optional, String, ForceNew) instance billing mode. Valid values: PREPAID (monthly subscription), POSTPAID_BY_HOUR (pay-as-you-go).
* `need_support_ipv6` - (Optional, Int, ForceNew) Whether to support IPv6 address access. Valid values: 1 (yes), 0 (no).
* `period` - (Optional, Int) Specify Prepaid period in month. Default `1`. Values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `voucher_ids` - (Optional, List: [`String`]) Specify Voucher Ids if `auto_voucher` was `1`, only support using 1 vouchers for now.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the postgresql instance.


## Import

postgresql readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_postgresql_readonly_instance.foo pgro-bcqx8b9a
```

