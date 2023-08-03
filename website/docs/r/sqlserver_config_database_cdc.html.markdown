---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_database_cdc"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_database_cdc"
description: |-
  Provides a resource to create a sqlserver config_database_cdc
---

# tencentcloud_sqlserver_config_database_cdc

Provides a resource to create a sqlserver config_database_cdc

## Example Usage

### Turn off database data change capture (CDC)

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

resource "tencentcloud_sqlserver_config_database_cdc" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  db_name     = tencentcloud_sqlserver_db.example.name
  modify_type = "disable"
}
```

### Enable Database Data Change Capture (CDC)

```hcl
resource "tencentcloud_sqlserver_config_database_cdc" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  db_name     = tencentcloud_sqlserver_db.example.name
  modify_type = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `db_name` - (Required, String) database name.
* `instance_id` - (Required, String) Instance ID.
* `modify_type` - (Required, String) Enable or disable CDC. Valid values: enable, disable.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver config_database_cdc can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_database_cdc.config_database_cdc config_database_cdc_id
```

