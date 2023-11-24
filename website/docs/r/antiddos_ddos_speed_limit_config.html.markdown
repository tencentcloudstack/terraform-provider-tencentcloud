---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_ddos_speed_limit_config"
sidebar_current: "docs-tencentcloud-resource-antiddos_ddos_speed_limit_config"
description: |-
  Provides a resource to create a antiddos ddos speed limit config
---

# tencentcloud_antiddos_ddos_speed_limit_config

Provides a resource to create a antiddos ddos speed limit config

## Example Usage

```hcl
resource "tencentcloud_antiddos_ddos_speed_limit_config" "ddos_speed_limit_config" {
  instance_id = "bgp-xxxxxx"
  ddos_speed_limit_config {
    mode = 1
    speed_values {
      type  = 1
      value = 1
    }
    speed_values {
      type  = 2
      value = 2
    }
    protocol_list = "ALL"
    dst_port_list = "8000"
  }
}
```

## Argument Reference

The following arguments are supported:

* `ddos_speed_limit_config` - (Required, List) Accessing speed limit configuration, the configuration ID cannot be empty when filling in parameters.
* `instance_id` - (Required, String) InstanceId.

The `ddos_speed_limit_config` object supports the following:

* `mode` - (Required, Int) Speed limit mode, value [1 (based on source IP speed limit) 2 (based on destination port speed limit)].
* `speed_values` - (Required, List) Speed limit values, each type of speed limit value can support up to 1; This field array has at least one speed limit value.
* `dst_port_list` - (Optional, String) List of port ranges, up to 8, multiple; Separate and indicate the range with -; This port range must be filled in; Fill in style 1:0-65535, style 2: 80; 443; 1000-2000.
* `dst_port_scopes` - (Optional, List) This field has been deprecated. Please fill in the new field DstPortList.
* `protocol_list` - (Optional, String) IP protocol numbers, values [ALL (all protocols) TCP (tcp protocol) UDP (udp protocol) SMP (smp protocol) 1; 2-100 (custom protocol number range, up to 8)] Note: When customizing the protocol number range, only the protocol number can be filled in, multiple ranges; Separation; When filling in ALL, no other agreements or agreements can be filled inNumber.

The `dst_port_scopes` object supports the following:

* `begin_port` - (Required, Int) Starting port, ranging from 1 to 65535.
* `end_port` - (Required, Int) end  port, ranging from 1 to 65535.

The `speed_values` object supports the following:

* `type` - (Required, Int) Speed limit value type, value [1 (packet rate pps) 2 (bandwidth bps)].
* `value` - (Required, Int) value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

antiddos ddos_speed_limit_config can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_ddos_speed_limit_config.ddos_speed_limit_config ${instanceId}#${configId}s
```

