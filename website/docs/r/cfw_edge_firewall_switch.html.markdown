---
subcategory: "Cfw"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_edge_firewall_switch"
sidebar_current: "docs-tencentcloud-resource-cfw_edge_firewall_switch"
description: |-
  Provides a resource to create a cfw edge_firewall_switch
---

# tencentcloud_cfw_edge_firewall_switch

Provides a resource to create a cfw edge_firewall_switch

## Example Usage

### If not set subnet_id

```hcl
data "tencentcloud_cfw_edge_fw_switches" "example" {}

resource "tencentcloud_cfw_edge_firewall_switch" "example" {
  public_ip   = data.tencentcloud_cfw_edge_fw_switches.example.data[0].public_ip
  switch_mode = 1
  enable      = 0
}
```

### If set subnet id

```hcl
data "tencentcloud_cfw_edge_fw_switches" "example" {}

resource "tencentcloud_cfw_edge_firewall_switch" "example" {
  public_ip   = data.tencentcloud_cfw_edge_fw_switches.example.data[0].public_ip
  subnet_id   = "subnet-id"
  switch_mode = 1
  enable      = 1
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Int) Switch, 0: off, 1: on.
* `public_ip` - (Required, String, ForceNew) Public Ip.
* `switch_mode` - (Required, Int) 0: bypass; 1: serial.
* `subnet_id` - (Optional, String) The first EIP switch in the vpc is turned on, and you need to specify a subnet to create a private connection. If `switch_mode` is 1 and `enable` is 1, this field is required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



