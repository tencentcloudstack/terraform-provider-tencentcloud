---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_port_acl_config"
sidebar_current: "docs-tencentcloud-resource-antiddos_port_acl_config"
description: |-
  Provides a resource to create a antiddos port acl config
---

# tencentcloud_antiddos_port_acl_config

Provides a resource to create a antiddos port acl config

## Example Usage

```hcl
resource "tencentcloud_antiddos_port_acl_config" "port_acl_config" {
  instance_id = "bgp-xxxxxx"
  acl_config {
    forward_protocol = "all"
    d_port_start     = 22
    d_port_end       = 23
    s_port_start     = 22
    s_port_end       = 23
    action           = "drop"
    priority         = 2

  }
}
```

## Argument Reference

The following arguments are supported:

* `acl_config` - (Required, List, ForceNew) Port ACL Policy.
* `instance_id` - (Required, String, ForceNew) InstanceIdList.

The `acl_config` object supports the following:

* `action` - (Required, String) Action, can take values: drop, transmit, forward.
* `d_port_end` - (Required, Int) end from port, with a range of 0~65535 values.
* `d_port_start` - (Required, Int) Starting from port, with a range of 0~65535 values.
* `forward_protocol` - (Required, String) Protocol type, can take TCP, udp, all values.
* `s_port_end` - (Required, Int) end from the source port, with a value range of 0~65535.
* `s_port_start` - (Required, Int) Starting from the source port, with a value range of 0~65535.
* `priority` - (Optional, Int) The policy priority, the smaller the number, the higher the level, and the higher the matching of the rule, with values ranging from 1 to 1000. Note: This field may return null, indicating that a valid value cannot be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

antiddos port_acl_config can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_port_acl_config.port_acl_config ${instanceId}#${configJson}
```

