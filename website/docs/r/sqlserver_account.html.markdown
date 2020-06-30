---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_account"
sidebar_current: "docs-tencentcloud-resource-sqlserver_account"
description: |-
  Use this resource to create SQL Server account
---

# tencentcloud_sqlserver_account

Use this resource to create SQL Server account

## Example Usage

```hcl
resource "tencentcloud_sqlserver_account" "foo" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-409mvdvv"
  subnet_id         = "subnet-nf9n81ps"
  engine_version    = "9.3.5"
  root_password     = "1qaA2k1wgvfa3ZZZ"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  storage           = 100
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Instance ID that the account belongs to.
* `name` - (Required) Name of the SQL Server account.
* `password` - (Required) Password of the SQL Server account.
* `is_admin` - (Optional) Indicate that the account is root account or not.
* `remark` - (Optional) Remark of the SQL Server account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the SQL Server account.
* `status` - Status of the SQL Server account. 1 for creating, 2 for running, 3 for modifying, 4 for resetting password, -1 for deleting.
* `update_time` - Last updated time of the SQL Server account.


## Import

sqlserver account can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_account.foo postgres-cda1iex1
```

