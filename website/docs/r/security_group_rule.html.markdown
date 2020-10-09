---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group_rule"
sidebar_current: "docs-tencentcloud-resource-security_group_rule"
description: |-
  Provides a resource to create security group rule.
---

# tencentcloud_security_group_rule

Provides a resource to create security group rule.

## Example Usage

Source is CIDR ip

```hcl
resource "tencentcloud_security_group" "sglab_1" {
  name        = "mysg_1"
  description = "favourite sg_1"
  project_id  = 0
}

resource "tencentcloud_security_group_rule" "sglab_1" {
  security_group_id = tencentcloud_security_group.sglab_1.id
  type              = "ingress"
  cidr_ip           = "10.0.0.0/16"
  ip_protocol       = "TCP"
  port_range        = "80"
  policy            = "ACCEPT"
  description       = "favourite sg rule_1"
}
```

Source is a security group id

```hcl
resource "tencentcloud_security_group" "sglab_2" {
  name        = "mysg_2"
  description = "favourite sg_2"
  project_id  = 0
}

resource "tencentcloud_security_group" "sglab_3" {
  name        = "mysg_3"
  description = "favourite sg_3"
  project_id  = 0
}

resource "tencentcloud_security_group_rule" "sglab_2" {
  security_group_id = tencentcloud_security_group.sglab_2.id
  type              = "ingress"
  ip_protocol       = "TCP"
  port_range        = "80"
  policy            = "ACCEPT"
  source_sgid       = tencentcloud_security_group.sglab_3.id
  description       = "favourite sg rule_2"
}
```

## Argument Reference

The following arguments are supported:

* `policy` - (Required, ForceNew) Rule policy of security group. Valid values: `ACCEPT` and `DROP`.
* `security_group_id` - (Required, ForceNew) ID of the security group to be queried.
* `type` - (Required, ForceNew) Type of the security group rule. Valid values: `ingress` and `egress`.
* `cidr_ip` - (Optional, ForceNew) An IP address network or segment, and conflict with `source_sgid`.
* `description` - (Optional, ForceNew) Description of the security group rule.
* `ip_protocol` - (Optional, ForceNew) Type of ip protocol. Valid values: `TCP`, `UDP` and `ICMP`. Default to all types protocol.
* `port_range` - (Optional, ForceNew) Range of the port. The available value can be one, multiple or one segment. E.g. `80`, `80,90` and `80-90`. Default to all ports.
* `source_sgid` - (Optional, ForceNew) ID of the nested security group, and conflict with `cidr_ip`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



