---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group_rule"
sidebar_current: "docs-tencentcloud-resource-security_group_rule"
description: |-
  Provide a resource to create security group rule.
---

# tencentcloud_security_group_rule

Provide a resource to create security group rule.

## Example Usage

```hcl
data "tencentcloud_security_group_rule" "sglab" {
    security_group_id = "sg-fh48e762"
    type              = "ingress"
    cidr_ip           = "10.0.0.0/16"
    ip_protocol       = "TCP"
    port_range        = "80"
    policy            = "ACCEPT"
    source_sgid       = "sg-fh48e762"
    description       = "favourite sg rule"
}
```

## Argument Reference

The following arguments are supported:

* `policy` - (Required, ForceNew) Rule policy of security group, the available value include 'ACCEPT' and 'DROP'.
* `security_group_id` - (Required, ForceNew) ID of the security group to be queried.
* `type` - (Required, ForceNew) Type of the security group rule, the available value include 'ingress' and 'egress'.
* `cidr_ip` - (Optional, ForceNew) An IP address network or segment, and can't exist at the same time as source_sgid.
* `description` - (Optional, ForceNew) Description of the security group rule.
* `ip_protocol` - (Optional, ForceNew) Type of ip protocol, the available value include 'TCP', 'UDP' and 'ICMP'. Default to all types.
* `port_range` - (Optional, ForceNew) Range of the port. The available value can be one, multiple or one segment, e.g '80', '80,90' and '80-90'. Default to all ports.
* `source_sgid` - (Optional, ForceNew) ID of the nested security group, and can't exist at the same time as cidr_ip.


