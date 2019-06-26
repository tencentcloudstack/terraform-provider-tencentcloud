---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group_rule"
sidebar_current: "docs-tencentcloud-resource-vpc-security-group-rule"
description: |-
  Provides an security group rule resource.
---

# tencentcloud_security_group_rule

Provides a security group rule resource. Represents a single `ingress` or `egress` group rule, which can be added to external Security Groups.

## Example Usage

Basic usage:

```hcl

resource "tencentcloud_security_group" "default" {
  name        = "${var.security_group_name}"
  description = "test security group rule"
}

resource "tencentcloud_security_group" "default2" {
  name        = "${var.security_group_name}"
  description = "anthor test security group rule"
}

resource "tencentcloud_security_group_rule" "http-in" {
  security_group_id = "${tencentcloud_security_group.default.id}"
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
}

resource "tencentcloud_security_group_rule" "ssh-in" {
  security_group_id = "${tencentcloud_security_group.default.id}"
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "22"
  policy            = "accept"
}

resource "tencentcloud_security_group_rule" "egress-drop" {
  security_group_id = "${tencentcloud_security_group.default.id}"
  type              = "egress"
  cidr_ip           = "10.2.3.0/24"
  ip_protocol       = "udp"
  port_range        = "3000-4000"
  policy            = "drop"
}

resource "tencentcloud_security_group_rule" "sourcesgid-in" {
  security_group_id = "${tencentcloud_security_group.default.id}"
  type              = "ingress"
  source_sgid       = "${tencentcloud_security_group.default2.id}"
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required, Forces new resource) The security group to apply this rule to.
* `type` - (Required, Forces new resource) The type of rule being created. Valid options are "ingress" (inbound) or "egress" (outbound).
* `cidr_ip` - (Optional, Forces new resource) can be IP, or CIDR block.
* `source_sgid` - (Optional, Forces new resource) The ID of a security group rule. Either `cidr_ip` or `source_sgid` must be specified, but it isn't supported simultaneously.
* `ip_protocol` - (Optional, Forces new resource) Support "UDP"、"TCP"、"ICMP", Not configured means all protocols.
* `port_range` - (Optional, Forces new resource) examples, Single port: "53"、Multiple ports: "80,8080,443"、Continuous port: "80-90", Not configured to represent all ports.
* `policy` - (Required, Forces new resource) Policy of rule, "accept" or "drop".

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the security group rule.
* `type` - The type of rule, "ingress" or "egress".
* `cidr_ip` - The source of rule, IP or CIDR block.
* `source_sgid` - The ID of a security group rule.
* `ip_protocol` – The protocol used.
* `port_range` – The port used.
* `policy` - The policy of rule, "accept" or "drop".
