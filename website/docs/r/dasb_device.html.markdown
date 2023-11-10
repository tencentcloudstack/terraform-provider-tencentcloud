---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_device"
sidebar_current: "docs-tencentcloud-resource-dasb_device"
description: |-
  Provides a resource to create a dasb device
---

# tencentcloud_dasb_device

Provides a resource to create a dasb device

## Example Usage

```hcl
resource "tencentcloud_dasb_device" "example" {
  os_name       = "Linux"
  ip            = "192.168.0.1"
  port          = 80
  name          = "tf_example"
  department_id = "1.2.3"
}
```

## Argument Reference

The following arguments are supported:

* `ip` - (Required, String) IP address.
* `os_name` - (Required, String) Operating system name, only Linux, Windows or MySQL.
* `port` - (Required, Int) Management port.
* `department_id` - (Optional, String) The department ID to which the device belongs.
* `ip_port_set` - (Optional, Set: [`String`]) Asset multi-node: fields ip and port.
* `name` - (Optional, String) Hostname, can be empty.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dasb device can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_device.example 17
```

