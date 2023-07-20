---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_instance"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_instance"
description: |-
  Provides a resource to create a monitor grafanaInstance
---

# tencentcloud_monitor_grafana_instance

Provides a resource to create a monitor grafanaInstance

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_monitor_grafana_instance" "foo" {
  instance_name         = "test-grafana"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_ids            = [tencentcloud_subnet.subnet.id]
  grafana_init_password = "1234567890"
  enable_internet       = false
  is_distroy            = true

  tags = {
    "createdBy" = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, String) Instance name.
* `enable_internet` - (Optional, Bool) Control whether grafana could be accessed by internet.
* `grafana_init_password` - (Optional, String) Grafana server admin password.
* `is_distroy` - (Optional, Bool) Whether to clean up completely, the default is false.
* `subnet_ids` - (Optional, Set: [`String`]) Subnet Id array.
* `tags` - (Optional, Map) Tag description list.
* `vpc_id` - (Optional, String) Vpc Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - Grafana instance id.
* `instance_status` - Grafana instance status, 1: Creating, 2: Running, 6: Stopped.
* `internal_url` - Grafana public address.
* `internet_url` - Grafana intranet address.
* `root_url` - Grafana external url which could be accessed by user.


## Import

monitor grafanaInstance can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_grafana_instance.foo grafanaInstance_id
```

