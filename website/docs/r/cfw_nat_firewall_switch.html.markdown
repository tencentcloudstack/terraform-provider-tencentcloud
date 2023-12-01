---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_nat_firewall_switch"
sidebar_current: "docs-tencentcloud-resource-cfw_nat_firewall_switch"
description: |-
  Provides a resource to create a cfw nat_firewall_switch
---

# tencentcloud_cfw_nat_firewall_switch

Provides a resource to create a cfw nat_firewall_switch

## Example Usage

### Turn off switch

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
}

resource "tencentcloud_cfw_nat_firewall_switch" "example" {
  nat_ins_id = data.tencentcloud_cfw_nat_fw_switches.example.id
  subnet_id  = data.tencentcloud_cfw_nat_fw_switches.example.data[0].subnet_id
  enable     = 0
}
```

### Or turn on switch

```hcl
data "tencentcloud_cfw_nat_fw_switches" "example" {
  nat_ins_id = "cfwnat-18d2ba18"
}

resource "tencentcloud_cfw_nat_firewall_switch" "example" {
  nat_ins_id = data.tencentcloud_cfw_nat_fw_switches.example.id
  subnet_id  = data.tencentcloud_cfw_nat_fw_switches.example.data[0].subnet_id
  enable     = 1
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Int) Switch, 0: off, 1: on.
* `nat_ins_id` - (Required, String, ForceNew) Firewall instance id.
* `subnet_id` - (Required, String, ForceNew) subnet id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cfw nat_firewall_switch can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_nat_firewall_switch.example cfwnat-18d2ba18#subnet-ef7wyymr
```

