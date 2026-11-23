---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group_rule"
sidebar_current: "docs-tencentcloud-resource-security_group_rule"
description: |-
  Provides a resource to create security group rule.
---

# tencentcloud_security_group_rule

Provides a resource to create security group rule.

~> **NOTE:** It is recommended to use `tencentcloud_security_group_rule_set` to manage multiple security group rules, because the index number of a single rule changes dynamically as other rules change.

## Example Usage

```hcl
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "security group example."
}

resource "tencentcloud_security_group_rule" "example" {
  security_group_id = tencentcloud_security_group.example.id
  type              = "ingress"
  cidr_ip           = "10.0.0.0/16"
  ip_protocol       = "TCP"
  port_range        = "80,8080"
  policy            = "ACCEPT"
  description       = "allow access to the web service."
}
```

## Argument Reference

The following arguments are supported:

* `policy` - (Required, String, ForceNew) Rule policy of security group. Valid values: `ACCEPT` and `DROP`.
* `security_group_id` - (Required, String, ForceNew) ID of the security group to be queried.
* `type` - (Required, String, ForceNew) Type of the security group rule. Valid values: `ingress` and `egress`.
* `address_template` - (Optional, List, ForceNew) ID of the address template, and confilicts with `source_sgid` and `cidr_ip`.
* `cidr_ip` - (Optional, String, ForceNew) An IP address network or segment, and conflict with `source_sgid` and `address_template`.
* `description` - (Optional, String, ForceNew) Description of the security group rule.
* `ip_protocol` - (Optional, String, ForceNew) Type of IP protocol. Valid values: `TCP`, `UDP`, `ICMP`, `ICMPv6` and `ALL`. Default to all types protocol, and conflicts with `protocol_template`.
* `policy_index` - (Optional, Int, ForceNew) The security group rule index number, the value of which dynamically changes as the security group rule changes.
* `port_range` - (Optional, String, ForceNew) Range of the port. The available value can be one, multiple or one segment. E.g. `80`, `80,90` and `80-90`. Default to all ports, and confilicts with `protocol_template`.
* `protocol_template` - (Optional, List, ForceNew) ID of the address template, and conflict with `ip_protocol`, `port_range`.
* `source_sgid` - (Optional, String, ForceNew) ID of the nested security group, and conflicts with `cidr_ip` and `address_template`.

The `address_template` object supports the following:

* `group_id` - (Optional, String, ForceNew) Address template group ID, conflicts with `template_id`.
* `template_id` - (Optional, String, ForceNew) Address template ID, conflicts with `group_id`.

The `protocol_template` object supports the following:

* `group_id` - (Optional, String, ForceNew) Address template group ID, conflicts with `template_id`.
* `template_id` - (Optional, String, ForceNew) Address template ID, conflicts with `group_id`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSecurityGroupPolicies](https://www.tencentcloud.com/document/api/215/15807)
* Read: [DescribeSecurityGroupPolicies](https://www.tencentcloud.com/document/api/215/15804)
* Delete:
  - [DeleteSecurityGroupPolicies](https://www.tencentcloud.com/document/api/215/15809)
  - [DescribeSecurityGroupPolicies](https://www.tencentcloud.com/document/api/215/15804)

