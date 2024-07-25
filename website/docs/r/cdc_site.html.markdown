---
subcategory: "CDC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdc_site"
sidebar_current: "docs-tencentcloud-resource-cdc_site"
description: |-
  Provides a resource to create a CDC site
---

# tencentcloud_cdc_site

Provides a resource to create a CDC site

## Example Usage

### Create a basic CDC site

```hcl
resource "tencentcloud_cdc_site" "example" {
  name         = "tf-example"
  country      = "China"
  province     = "Guangdong Province"
  city         = "Guangzhou"
  address_line = "Tencent Building"
  description  = "desc."
}
```

### Create a complete CDC site

```hcl
resource "tencentcloud_cdc_site" "example" {
  name                  = "tf-example"
  country               = "China"
  province              = "Guangdong Province"
  city                  = "Guangzhou"
  address_line          = "Shenzhen Tencent Building"
  optional_address_line = "Shenzhen Tencent Building of Binhai"
  description           = "desc."
  fiber_type            = "MM"
  optical_standard      = "MM"
  power_connectors      = "380VAC3P"
  power_feed_drop       = "DOWN"
  max_weight            = 100
  power_draw_kva        = 10
  uplink_speed_gbps     = 10
  uplink_count          = 2
  condition_requirement = true
  dimension_requirement = true
  redundant_networking  = true
  need_help             = true
  redundant_power       = true
  breaker_requirement   = true
}
```

## Argument Reference

The following arguments are supported:

* `address_line` - (Required, String) Site Detail Address.
* `city` - (Required, String) Site City.
* `country` - (Required, String) Site Country.
* `name` - (Required, String) Site Name.
* `province` - (Required, String) Site Province.
* `breaker_requirement` - (Optional, Bool) Whether there is an upstream circuit breaker.
* `condition_requirement` - (Optional, Bool) Whether the following environmental conditions are met: n1. There are no material requirements or the acceptance standard on site that will affect the delivery and installation of the CDC device. n2. The following conditions are met for finalized rack positions: Temperature ranges from 41 to 104 degrees F (5 to 40 degrees C). Humidity ranges from 10 degrees F (-12 degrees C) to 70 degrees F (21 degrees C) and relative humidity ranges from 8% RH to 80% RH. Air flows from front to back at the rack position and there is sufficient air in CFM (cubic feet per minute). The air quantity in CFM must be 145.8 times the power consumption (in KVA) of CDC.
* `description` - (Optional, String) Site Description.
* `dimension_requirement` - (Optional, Bool) Whether the following dimension conditions are met: Your loading dock can accommodate one rack container (H x W x D = 94 x 54 x 48). You can provide a clear route from the delivery point of your rack (H x W x D = 80 x 24 x 48) to its final installation location. You should consider platforms, corridors, doors, turns, ramps, freight elevators as well as other access restrictions when measuring the depth. There shall be a 48 or greater front clearance and a 24 or greater rear clearance where the CDC is finally installed.
* `fiber_type` - (Optional, String) Site Fiber Type. Using optical fiber type to connect the CDC device to the network SM(Single-Mode) or MM(Multi-Mode) fibers are available.
* `max_weight` - (Optional, Int) Site Max Weight capacity (KG).
* `need_help` - (Optional, Bool) Whether you need help from Tencent Cloud for rack installation.
* `optical_standard` - (Optional, String) Site Optical Standard. Optical standard used to connect the CDC device to the network This field depends on the uplink speed, optical fiber type, and distance to upstream equipment. Allow value: `SM`, `MM`.
* `optional_address_line` - (Optional, String) Detailed address of the site area (to be added).
* `power_connectors` - (Optional, String) Site Power Connectors. Example: 380VAC3P.
* `power_draw_kva` - (Optional, Int) Site Power DrawKva (KW).
* `power_feed_drop` - (Optional, String) Site Power Feed Drop. Whether power is supplied from above or below the rack. Allow value: `UP`, `DOWN`.
* `redundant_networking` - (Optional, Bool) Whether redundant upstream equipment (switch or router) is provided so that both network devices can be connected to the network.
* `redundant_power` - (Optional, Bool) Whether there is power redundancy.
* `uplink_count` - (Optional, Int) Number of uplinks used by each CDC device (2 devices per rack) when connected to the network.
* `uplink_speed_gbps` - (Optional, Int) Uplink speed from the network to Tencent Cloud Region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CDC site can be imported using the id, e.g.

```
terraform import tencentcloud_cdc_site.example site-43qcf1ag
```

