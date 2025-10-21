---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group_lite_rule"
sidebar_current: "docs-tencentcloud-resource-security_group_lite_rule"
description: |-
  Provide a resource to create security group some lite rules quickly.
---

# tencentcloud_security_group_lite_rule

Provide a resource to create security group some lite rules quickly.

-> **NOTE:** This resource has been deprecated in Terraform TencentCloud provider version 1.81.90. Please use 'tencentcloud_security_group_rule_set' instead.

## Example Usage

```hcl
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
}

resource "tencentcloud_security_group_lite_rule" "foo" {
  security_group_id = tencentcloud_security_group.foo.id

  ingress = [
    "ACCEPT#192.168.1.0/24#80#TCP",
    "DROP#8.8.8.8#80,90#UDP",
    "ACCEPT#0.0.0.0/0#80-90#TCP",
    "ACCEPT#sg-7ixn3foj#80-90#TCP",
    "ACCEPT#ipm-epjq5kn0#80-90#TCP",
    "ACCEPT#ipmg-3loavam6#80-90#TCP",
    "ACCEPT#0.0.0.0/0##ppm-xxxxxxxx"
    "ACCEPT#0.0.0.0/0##ppmg-xxxxxxxx"
  ]

  egress = [
    "ACCEPT#192.168.0.0/16#ALL#TCP",
    "ACCEPT#10.0.0.0/8#ALL#ICMP",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required, String, ForceNew) ID of the security group.
* `egress` - (Optional, List: [`String`]) Egress rules set. A rule must match the following format: [action]#[source]#[port]#[protocol]. The available value of 'action' is `ACCEPT` and `DROP`. The 'source' can be an IP address network, segment, security group ID and Address Template ID. The 'port' valid format is `80`, `80,443`, `80-90` or `ALL`. The available value of 'protocol' is `TCP`, `UDP`, `ICMP`, `ALL` and `ppm(g?)-xxxxxxxx`. When 'protocol' is `ICMP` or `ALL`, the 'port' must be `ALL`.
* `ingress` - (Optional, List: [`String`]) Ingress rules set. A rule must match the following format: [action]#[source]#[port]#[protocol]. The available value of 'action' is `ACCEPT` and `DROP`. The 'source' can be an IP address network, segment, security group ID and Address Template ID. The 'port' valid format is `80`, `80,443`, `80-90` or `ALL`. The available value of 'protocol' is `TCP`, `UDP`, `ICMP`, `ALL` and `ppm(g?)-xxxxxxxx`. When 'protocol' is `ICMP` or `ALL`, the 'port' must be `ALL`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Security group lite rule can be imported using the id, e.g.

```
  $ terraform import tencentcloud_security_group_lite_rule.foo sg-ey3wmiz1
```

