---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_start_xevent"
sidebar_current: "docs-tencentcloud-resource-sqlserver_start_xevent"
description: |-
  Provides a resource to create a sqlserver start_xevent
---

# tencentcloud_sqlserver_start_xevent

Provides a resource to create a sqlserver start_xevent

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

resource "tencentcloud_sqlserver_start_xevent" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  event_config {
    event_type = "slow"
    threshold  = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `event_config` - (Required, List, ForceNew) Whether to start or stop an extended event.
* `instance_id` - (Required, String, ForceNew) Instance ID.

The `event_config` object supports the following:

* `event_type` - (Required, String) Event type. Valid values: slow (set threshold for slow SQL ), blocked (set threshold for the blocking and deadlock).
* `threshold` - (Required, Int) Threshold in milliseconds. Valid values: 0(disable), non-zero (enable).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



