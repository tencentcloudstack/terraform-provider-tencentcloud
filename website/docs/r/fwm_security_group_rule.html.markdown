---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_security_group_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_security_group_rule"
description: |-
  Provides a resource to create a FWM security group rule.
---

# tencentcloud_fwm_security_group_rule

Provides a resource to create a FWM security group rule.

## Example Usage

```hcl
resource "tencentcloud_fwm_security_group_rule" "example" {
  group_id           = "fwmrg-ghlu11dv"
  ip_version         = "ipv4"
  source_type        = "net"
  source_content     = "0.0.0.0/0"
  dest_type          = "net"
  dest_content       = "0.0.0.0/0"
  protocol           = "ANY"
  port               = "-1/-1"
  rule_action        = "drop"
  scope              = "SG"
  protocol_port_type = 0
  description        = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `dest_content` - (Required, String) Destination content, such as IP, CIDR or instance ID.
* `dest_type` - (Required, String, ForceNew) Destination type: net, vpc, subnet, cvm, clb, eni, cdb, template, tag, region, resourcegroup, addressbook, instance.
* `group_id` - (Required, String, ForceNew) Rule group ID, the parent rule group (Product=enterprise_sg) to which the rule belongs.
* `ip_version` - (Required, String, ForceNew) IP version: ipv4, ipv6.
* `port` - (Required, String) Port, such as 80, 1-65535 or ALL.
* `protocol` - (Required, String) Protocol: TCP, UDP, ICMP, ANY. ANY means all protocols.
* `rule_action` - (Required, String) Rule action: accept, drop, log.
* `source_content` - (Required, String) Source content, such as IP, CIDR or instance ID.
* `source_type` - (Required, String, ForceNew) Source type: net, vpc, subnet, cvm, clb, eni, cdb, template, tag, region, resourcegroup, addressbook, instance.
* `belong_member_id` - (Optional, String) ID of the member account that the rule belongs to.
* `description` - (Optional, String) Rule description.
* `protocol_port_type` - (Optional, Int) Protocol port type.
* `scope` - (Optional, String, ForceNew) Rule effective scope: SG, LH.
* `service_template_id` - (Optional, String) Service template ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.
* `sequence` - Priority (read-only, same source as order_index).
* `source_id` - Source ID.
* `target_id` - Destination ID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSecurityGroupRule](https://www.tencentcloud.com/search/?wd=CreateSecurityGroupRule&type=document&lang=en)
* Read: [DescribeSecurityGroupRules](https://www.tencentcloud.com/search/?wd=DescribeSecurityGroupRules&type=document&lang=en)
* Update: [ModifySecurityGroupRule](https://www.tencentcloud.com/search/?wd=ModifySecurityGroupRule&type=document&lang=en)
* Delete: [DeleteSecurityGroupRule](https://www.tencentcloud.com/search/?wd=DeleteSecurityGroupRule&type=document&lang=en)

## Import

FWM security group rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_security_group_rule.example fwmrg-ghlu11dv#a5f39ae7-2a81-4735-9027-5d6f83d814a5
```

