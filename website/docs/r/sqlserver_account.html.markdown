---
subcategory: "SQLServer"
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
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.1.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_sqlserver_instance" "example" {
  name              = "tf-example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.1.name
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  project_id        = 0
  memory            = 4
  storage           = 100
}

resource "tencentcloud_sqlserver_account" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name        = "tf_example_account"
  password    = "Qwer@234"
  remark      = "test-remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID that the account belongs to.
* `name` - (Required, String) Name of the SQL Server account.
* `password` - (Required, String) Password of the SQL Server account.
* `is_admin` - (Optional, Bool) Indicate that the account is root account or not.
* `remark` - (Optional, String) Remark of the SQL Server account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the SQL Server account.
* `status` - Status of the SQL Server account. Valid values: 1, 2, 3, 4. 1 for creating, 2 for running, 3 for modifying, 4 for resetting password, -1 for deleting.
* `update_time` - Last updated time of the SQL Server account.


## Import

SQL Server account can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_account.foo mssql-3cdq7kx5#tf_sqlserver_account
```

