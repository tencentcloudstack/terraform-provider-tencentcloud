---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_account_db_attachment"
sidebar_current: "docs-tencentcloud-resource-sqlserver_account_db_attachment"
description: |-
  Use this resource to create SQL Server account DB attachment
---

# tencentcloud_sqlserver_account_db_attachment

Use this resource to create SQL Server account DB attachment

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

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name        = "tf_example_db"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_account" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name        = "tf_example_account"
  password    = "Qwer@234"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_account_db_attachment" "example" {
  instance_id  = tencentcloud_sqlserver_instance.example.id
  account_name = tencentcloud_sqlserver_account.example.name
  db_name      = tencentcloud_sqlserver_db.example.name
  privilege    = "ReadWrite"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, String, ForceNew) SQL Server account name.
* `db_name` - (Required, String, ForceNew) SQL Server DB name.
* `instance_id` - (Required, String, ForceNew) SQL Server instance ID that the account belongs to.
* `privilege` - (Required, String) Privilege of the account on DB. Valid values: `ReadOnly`, `ReadWrite`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SQL Server account DB attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_account_db_attachment.foo mssql-3cdq7kx5#tf_sqlserver_account#test111
```

