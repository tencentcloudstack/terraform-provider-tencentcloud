---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group_rule_set"
sidebar_current: "docs-tencentcloud-resource-security_group_rule_set"
description: |-
  Provides a resource to create security group rule. This resource is similar with tencentcloud_security_group_lite_rule, rules can be ordered and configure descriptions.
---

# tencentcloud_security_group_rule_set

Provides a resource to create security group rule. This resource is similar with tencentcloud_security_group_lite_rule, rules can be ordered and configure descriptions.

~> **NOTE:** This resource must exclusive in one security group, do not declare additional rule resources of this security group elsewhere.

## Example Usage

```hcl
resource "tencentcloud_security_group" "base" {
  name        = "test-set-sg"
  description = "Testing Rule Set Security"
}

resource "tencentcloud_security_group" "relative" {
  name        = "for-relative"
  description = "Used for attach security policy"
}

resource "tencentcloud_address_template" "foo" {
  name      = "test-set-aTemp"
  addresses = ["10.0.0.1", "10.0.1.0/24", "10.0.0.1-10.0.0.100"]
}

resource "tencentcloud_address_template_group" "foo" {
  name         = "test-set-atg"
  template_ids = [tencentcloud_address_template.foo.id]
}

resource "tencentcloud_security_group_rule_set" "base" {
  security_group_id = tencentcloud_security_group.base.id

  ingress {
    action      = "ACCEPT"
    cidr_block  = "10.0.0.0/22"
    protocol    = "TCP"
    port        = "80-90"
    description = "A:Allow Ips and 80-90"
  }

  ingress {
    action      = "ACCEPT"
    cidr_block  = "10.0.2.1"
    protocol    = "UDP"
    port        = "8080"
    description = "B:Allow UDP 8080"
  }

  ingress {
    action      = "ACCEPT"
    cidr_block  = "10.0.2.1"
    protocol    = "UDP"
    port        = "8080"
    description = "C:Allow UDP 8080"
  }

  ingress {
    action      = "ACCEPT"
    cidr_block  = "172.18.1.2"
    protocol    = "ALL"
    port        = "ALL"
    description = "D:Allow ALL"
  }

  ingress {
    action             = "DROP"
    protocol           = "TCP"
    port               = "80"
    source_security_id = tencentcloud_security_group.relative.id
    description        = "E:Block relative"
  }

  egress {
    action      = "DROP"
    cidr_block  = "10.0.0.0/16"
    protocol    = "ICMP"
    description = "A:Block ping3"
  }

  egress {
    action              = "DROP"
    address_template_id = tencentcloud_address_template.foo.id
    description         = "B:Allow template"
  }

  egress {
    action                 = "DROP"
    address_template_group = tencentcloud_address_template_group.foo.id
    description            = "C:DROP template group"
  }
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required, String, ForceNew) ID of the security group to be queried.
* `egress` - (Optional, List) List of egress rule. NOTE: this block is ordered, the first rule has the highest priority.
* `ingress` - (Optional, List) List of ingress rule. NOTE: this block is ordered, the first rule has the highest priority.

The `egress` object supports the following:

* `action` - (Required, String) Rule policy of security group. Valid values: `ACCEPT` and `DROP`.
* `address_template_group` - (Optional, String) Specify Group ID of Address template like `ipmg-xxxxxxxx`, conflict with `source_security_id` and `cidr_block`.
* `address_template_id` - (Optional, String) Specify Address template ID like `ipm-xxxxxxxx`, conflict with `source_security_id` and `cidr_block`.
* `cidr_block` - (Optional, String) An IP address network or CIDR segment. NOTE: `cidr_block`, `ipv6_cidr_block`, `source_security_id` and `address_template_*` are exclusive and cannot be set in the same time.
* `description` - (Optional, String) Description of the security group rule.
* `ipv6_cidr_block` - (Optional, String) An IPV6 address network or CIDR segment, and conflict with `source_security_id` and `address_template_*`.
* `port` - (Optional, String) Range of the port. The available value can be one, multiple or one segment. E.g. `80`, `80,90` and `80-90`. Default to all ports, and conflicts with `service_template_*`.
* `protocol` - (Optional, String) Type of IP protocol. Valid values: `TCP`, `UDP` and `ICMP`. Default to all types protocol, and conflicts with `service_template_*`.
* `service_template_group` - (Optional, String) Specify Group ID of Protocol template ID like `ppmg-xxxxxxxx`, conflict with `cidr_block` and `port`.
* `service_template_id` - (Optional, String) Specify Protocol template ID like `ppm-xxxxxxxx`, conflict with `cidr_block` and `port`.
* `source_security_id` - (Optional, String) ID of the nested security group, and conflicts with `cidr_block` and `address_template_*`.

The `ingress` object supports the following:

* `action` - (Required, String) Rule policy of security group. Valid values: `ACCEPT` and `DROP`.
* `address_template_group` - (Optional, String) Specify Group ID of Address template like `ipmg-xxxxxxxx`, conflict with `source_security_id` and `cidr_block`.
* `address_template_id` - (Optional, String) Specify Address template ID like `ipm-xxxxxxxx`, conflict with `source_security_id` and `cidr_block`.
* `cidr_block` - (Optional, String) An IP address network or CIDR segment. NOTE: `cidr_block`, `ipv6_cidr_block`, `source_security_id` and `address_template_*` are exclusive and cannot be set in the same time.
* `description` - (Optional, String) Description of the security group rule.
* `ipv6_cidr_block` - (Optional, String) An IPV6 address network or CIDR segment, and conflict with `source_security_id` and `address_template_*`.
* `port` - (Optional, String) Range of the port. The available value can be one, multiple or one segment. E.g. `80`, `80,90` and `80-90`. Default to all ports, and conflicts with `service_template_*`.
* `protocol` - (Optional, String) Type of IP protocol. Valid values: `TCP`, `UDP` and `ICMP`. Default to all types protocol, and conflicts with `service_template_*`.
* `service_template_group` - (Optional, String) Specify Group ID of Protocol template ID like `ppmg-xxxxxxxx`, conflict with `cidr_block` and `port`.
* `service_template_id` - (Optional, String) Specify Protocol template ID like `ppm-xxxxxxxx`, conflict with `cidr_block` and `port`.
* `source_security_id` - (Optional, String) ID of the nested security group, and conflicts with `cidr_block` and `address_template_*`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `version` - Security policies version, auto increment for every update.


## Import

Resource tencentcloud_security_group_rule_set can be imported by passing security grou id:

```
terraform import tencentcloud_security_group_rule_set.sglab_1 sg-xxxxxxxx
```

