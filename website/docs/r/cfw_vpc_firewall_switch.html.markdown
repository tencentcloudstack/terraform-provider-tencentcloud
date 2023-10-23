---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_vpc_firewall_switch"
sidebar_current: "docs-tencentcloud-resource-cfw_vpc_firewall_switch"
description: |-
  Provides a resource to create a cfw vpc_firewall_switch
---

# tencentcloud_cfw_vpc_firewall_switch

Provides a resource to create a cfw vpc_firewall_switch

## Example Usage

### Turn off switch

```hcl
data "tencentcloud_cfw_vpc_fw_switches" "example" {
  vpc_ins_id = "cfwg-c8c2de41"
}

resource "tencentcloud_cfw_vpc_firewall_switch" "example" {
  vpc_ins_id = data.tencentcloud_cfw_vpc_fw_switches.example.id
  switch_id  = data.tencentcloud_cfw_vpc_fw_switches.example.switch_list[0].switch_id
  enable     = 0
}
```

### Or turn on switch

```hcl
data "tencentcloud_cfw_vpc_fw_switches" "example" {
  vpc_ins_id = "cfwg-c8c2de41"
}

resource "tencentcloud_cfw_vpc_firewall_switch" "example" {
  vpc_ins_id = data.tencentcloud_cfw_vpc_fw_switches.example.id
  switch_id  = data.tencentcloud_cfw_vpc_fw_switches.example.switch_list[0].switch_id
  enable     = 1
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Int) Turn the switch on or off. 0: turn off the switch; 1: Turn on the switch.
* `switch_id` - (Required, String, ForceNew) Firewall switch ID.
* `vpc_ins_id` - (Required, String, ForceNew) Firewall instance id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cfw vpc_firewall_switch can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_vpc_firewall_switch.example cfwg-c8c2de41#cfws-f2c63ded84
```

