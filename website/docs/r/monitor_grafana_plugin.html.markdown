---
subcategory: "TencentCloud Managed Service for Grafana(TCMG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_grafana_plugin"
sidebar_current: "docs-tencentcloud-resource-monitor_grafana_plugin"
description: |-
  Provides a resource to create a monitor grafanaPlugin
---

# tencentcloud_monitor_grafana_plugin

Provides a resource to create a monitor grafanaPlugin

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

  tags = {
    "createdBy" = "test"
  }
}

resource "tencentcloud_monitor_grafana_plugin" "grafanaPlugin" {
  instance_id = tencentcloud_monitor_grafana_instance.foo.id
  plugin_id   = "grafana-piechart-panel"
  version     = "1.6.2"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Grafana instance id.
* `plugin_id` - (Required, String) Plugin id.
* `version` - (Optional, String) Plugin version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor grafanaPlugin can be imported using the instance_id#plugin_id, e.g.
```
$ terraform import tencentcloud_monitor_grafana_plugin.grafanaPlugin grafana-50nj6v00#grafana-piechart-panel
```

