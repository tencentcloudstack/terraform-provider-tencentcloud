---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_db"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_db"
description: |-
  Use this data source to query DB resources for the specific SQLServer instance.
---

# tencentcloud_sqlserver_db

Use this data source to query DB resources for the specific SQLServer instance.

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-2"
}

resource "tencentcloud_vpc" "foo" {
  name       = "example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  name              = "example"
  availability_zone = var.availability_zone
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.0.0/24"
  is_multicast      = false
}

resource "tencentcloud_sqlserver_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.foo.id
  subnet_id         = tencentcloud_subnet.foo.id
  engine_version    = "2008R2"
  project_id        = 0
  memory            = 2
  storage           = 10
}

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name        = "example"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

data "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
  name        = tencentcloud_sqlserver_db.example.name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) SQLServer instance ID which DB belongs to.
* `name` - (Required) Name of DB.
* `result_output_file` - (Optional) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `charset` - Character set DB uses, could be `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `Chinese_PRC_BIN`, `Chinese_Taiwan_Stroke_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`, and `SQL_Latin1_General_CP1_CS_AS`. Default value is `Chinese_PRC_CI_AS`.
* `create_time` - Database creation time.
* `remark` - Remark of the DB.
* `status` - Database status. Valid values are `creating`, `running`, `modifying`, `dropping`.


