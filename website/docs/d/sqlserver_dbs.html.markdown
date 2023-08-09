---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_dbs"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_dbs"
description: |-
  Use this data source to query DB resources for the specific SQL Server instance.
---

# tencentcloud_sqlserver_dbs

Use this data source to query DB resources for the specific SQL Server instance.

## Example Usage

```hcl
data "tencentcloud_sqlserver_dbs" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
}

data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_db"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) SQL Server instance ID which DB belongs to.
* `result_output_file` - (Optional, String) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `db_list` - A list of dbs belong to the specific instance. Each element contains the following attributes:
  * `charset` - Character set DB uses, could be `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `Chinese_PRC_BIN`, `Chinese_Taiwan_Stroke_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`, and `SQL_Latin1_General_CP1_CS_AS`.
  * `create_time` - Database creation time.
  * `instance_id` - SQL Server instance ID which DB belongs to.
  * `name` - Name of DB.
  * `remark` - Remark of the DB.
  * `status` - Database status. Valid values are `creating`, `running`, `modifying`, `dropping`.


