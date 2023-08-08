---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_database_tde"
sidebar_current: "docs-tencentcloud-resource-sqlserver_database_tde"
description: |-
  Provides a resource to create a sqlserver database_tde
---

# tencentcloud_sqlserver_database_tde

Provides a resource to create a sqlserver database_tde

## Example Usage

### Open database tde encryption

```hcl
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

resource "tencentcloud_sqlserver_database_tde" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  db_names    = [tencentcloud_sqlserver_db.example.name]
  encryption  = "enable"
}
```

### Close database tde encryption

```hcl
resource "tencentcloud_sqlserver_database_tde" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  db_names    = [tencentcloud_sqlserver_db.example.name]
  encryption  = "disable"
}
```

## Argument Reference

The following arguments are supported:

* `db_names` - (Required, Set: [`String`]) Database name list.
* `encryption` - (Required, String) `enable` - enable encryption, `disable` - disable encryption.
* `instance_id` - (Required, String) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver database_tde can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_database_tde.database_tde database_tde_id
```

