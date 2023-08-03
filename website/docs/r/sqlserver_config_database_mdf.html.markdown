---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_database_mdf"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_database_mdf"
description: |-
  Provides a resource to create a sqlserver config_database_mdf
---

# tencentcloud_sqlserver_config_database_mdf

Provides a resource to create a sqlserver config_database_mdf

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
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_sqlserver_instance" "example" {
  name              = "tf-example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  project_id        = 0
  memory            = 16
  storage           = 40
}

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name        = "tf_example_db"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_config_database_mdf" "example" {
  db_name     = tencentcloud_sqlserver_db.example.name
  instance_id = tencentcloud_sqlserver_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `db_name` - (Required, String) Array of database names.
* `instance_id` - (Required, String) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver config_database_mdf can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_database_mdf.config_database_mdf config_database_mdf_id
```

