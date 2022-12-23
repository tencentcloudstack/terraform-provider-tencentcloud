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
resource "tencentcloud_security_group" "sglab_1" {
  name        = "mysg_1"
  description = "favourite sg_1"
}

resource "tencentcloud_security_group_rule_set" "sglab_1" {
  security_group_id = tencentcloud_security_group.sglab_1.id
  ingress {
    cidr_block  = "10.0.0.0/16" # Accept IP or CIDR
    protocol    = "TCP"         # Default is ALL
    port        = "80"          # Accept port e.g. 80 or PortRange e.g. 8080-8089
    action      = "ACCEPT"
    description = "favourite sg rule_1"
  }
  ingress {
    protocol           = "TCP"
    port               = "80"
    action             = "ACCEPT"
    source_security_id = tencentcloud_security_group.sglab_3.id
    description        = "favourite sg rule_2"
  }

  egress {
    action              = "ACCEPT"
    address_template_id = "ipm-xxxxxxxx" # Support address template (group)
    description         = "Allow address template"
  }
  egress {
    action                 = "ACCEPT"
    service_template_group = "ppmg-xxxxxxxx" # Support protocol template (group)
    description            = "Allow protocol template"
  }
  egress {
    cidr_block  = "10.0.0.0/16"
    protocol    = "TCP"
    port        = "80"
    action      = "DROP"
    description = "favourite sg egress rule"
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

