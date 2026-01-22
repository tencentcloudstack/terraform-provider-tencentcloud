---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_parameters"
sidebar_current: "docs-tencentcloud-resource-postgresql_parameters"
description: |-
  Use this resource to create postgresql parameter.
---

# tencentcloud_postgresql_parameters

Use this resource to create postgresql parameter.

## Example Usage

```hcl
variable "default_az" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "gz3" {
  availability_zone = var.default_az
  is_default        = true
}

locals {
  vpc_id    = data.tencentcloud_vpc_subnets.gz3.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.subnet_id
}

data "tencentcloud_availability_zones_by_product" "zone" {
  product = "postgres"
}

resource "tencentcloud_postgresql_instance" "test" {
  name              = "tf_postsql_postpaid"
  availability_zone = var.default_az
  charge_type       = "POSTPAID_BY_HOUR"
  period            = 1
  vpc_id            = local.vpc_id
  subnet_id         = local.subnet_id
  engine_version    = "13.3"
  root_password     = "t1qaA2k1wgvfa3?ZZZ"
  security_groups   = ["sg-5275dorp"]
  charset           = "LATIN1"
  project_id        = 0
  memory            = 2
  storage           = 20
}
resource "tencentcloud_postgresql_parameters" "postgresql_parameters" {
  db_instance_id = tencentcloud_postgresql_instance.test.id
  param_list {
    expected_value = "off"
    name           = "check_function_bodies"
  }
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance ID.
* `param_list` - (Required, List) Parameters to be modified and expected values.

The `param_list` object supports the following:

* `expected_value` - (Required, String) The new value to which the parameter will be modified. When this parameter is used as an input parameter, its value must be a string, such as `0.1` (decimal), `1000` (integer), and `replica` (enum).
* `name` - (Required, String) Parameter name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgresql parameters can be imported, e.g.

```
$ terraform import tencentcloud_postgresql_parameters.example pgrogrp-lckioi2a
```

