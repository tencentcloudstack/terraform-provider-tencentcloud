---
subcategory: "PostgreSQL"
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

* `db_version` - (Required, ForceNew) PostgreSQL kernel version, which must be the same as that of the primary instance.
* `master_db_instance_id` - (Required, ForceNew) ID of the primary instance to which the read-only replica belongs.
* `memory` - (Required) Memory size(in GB). Allowed value must be larger than `memory` that data source `tencentcloud_postgresql_specinfos` provides.
* `name` - (Required) Instance name.
* `project_id` - (Required) Project ID.
* `security_groups_ids` - (Required) ID of security group.
* `storage` - (Required) Instance storage capacity in GB.
* `subnet_id` - (Required) VPC subnet ID.
* `vpc_id` - (Required, ForceNew) VPC ID.
* `zone` - (Required, ForceNew) Availability zone ID, which can be obtained through the Zone field in the returned value of the DescribeZones API.
* `auto_renew_flag` - (Optional, ForceNew) Renewal flag. Valid values: 0 (manual renewal), 1 (auto-renewal). Default value: 0.
* `instance_charge_type` - (Optional, ForceNew) instance billing mode. Valid values: PREPAID (monthly subscription), POSTPAID_BY_HOUR (pay-as-you-go).
* `need_support_ipv6` - (Optional, ForceNew) Whether to support IPv6 address access. Valid values: 1 (yes), 0 (no).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the postgresql instance.


## Import

postgresql readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_postgresql_readonly_instance.foo pgro-bcqx8b9a
```

