---
subcategory: "Cloud Firewall(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_security_group_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_security_group_rule"
description: |-
  Provides a resource to create a FWM security group rule.
---

# tencentcloud_fwm_security_group_rule

Provides a resource to create a FWM security group rule.

~> **NOTE:** The `tencentcloud_fwm_security_group_rule` resource belongs to a rule group (Product=enterprise_sg) which must be created in advance by `tencentcloud_fwm_rule_group`. This resource does not create the rule group.

## Example Usage

```hcl
resource "tencentcloud_fwm_security_group_rule" "example" {
  group_id       = "group-id-example"
  ip_version     = "ipv4"
  order_index    = 1
  source_type    = "net"
  source_content = "1.1.1.1"
  dest_type      = "net"
  dest_content   = "2.2.2.2"
  protocol       = "TCP"
  port           = "80"
  rule_action    = "accept"
  scope          = "SG"
  description    = "tf example"
}
```

## Argument Reference

The following arguments are supported:

* `dest_content` - (Required, String) Destination content, such as IP/CIDR or instance id.
* `dest_type` - (Required, String, ForceNew) Destination type, optional values: net, vpc, subnet, cvm, clb, eni, cdb, template, tag, region, resourcegroup, addressbook, instance.
* `group_id` - (Required, String) Rule group ID, the parent rule group (Product=enterprise_sg) to which the rule belongs.
* `ip_version` - (Required, String, ForceNew) IP version, optional values: ipv4, ipv6.
* `order_index` - (Required, Int) Rule priority order index, global single sequence within the group, >=1.
* `port` - (Required, String) Port, such as 80 or ALL.
* `protocol` - (Required, String) Protocol, optional values: TCP, ANY, ICMP, UDP. ANY means all.
* `rule_action` - (Required, String) Rule action, optional values: accept, drop, log.
* `source_content` - (Required, String) Source content, such as IP/CIDR or instance id.
* `source_type` - (Required, String, ForceNew) Source type, optional values: net, vpc, subnet, cvm, clb, eni, cdb, template, tag, region, resourcegroup, addressbook, instance.
* `belong_member_id` - (Optional, String) Member id to which the rule belongs.
* `description` - (Optional, String) Rule description.
* `protocol_port_type` - (Optional, Int) Protocol port type.
* `scope` - (Optional, String, ForceNew) Rule effective scope, optional values: SG, LH.
* `service_template_id` - (Optional, String) Service template id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `all_total_count` - All total count (read-only display).
* `cidr` - Destination cidr (read-only display).
* `detail` - Rule detail (read-only display).
* `dns_parse_count` - Dns parse count (read-only display).
  * `invalid_count` - Invalid ip count.
  * `valid_count` - Valid ip count.
* `id` - Internal primary key id (read-only display).
* `instance_name` - Destination instance name (read-only display).
* `is_new` - Whether the rule has been changed recently (read-only display).
* `parameter_name` - Destination parameter name (read-only display).
* `private_ip` - Destination private ip (read-only display).
* `protocol_port_name` - Protocol port name (read-only display).
* `public_ip` - Destination public ip (read-only display).
* `region` - Region (read-only display).
* `rule_id` - Rule id.
* `sou_cidr` - Source cidr (read-only display).
* `sou_instance_name` - Source instance name (read-only display).
* `sou_parameter_name` - Source parameter name (read-only display).
* `sou_private_ip` - Source private ip (read-only display).
* `sou_public_ip` - Source public ip (read-only display).
* `source_id` - Source id (read-only display).
* `strategy` - Rule strategy int value (read-only display).
* `target_id` - Target id (read-only display).

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSecurityGroupRule](https://www.tencentcloud.com/search/?wd=CreateSecurityGroupRule&type=document&lang=en)
* Read: [DescribeSecurityGroupRules](https://www.tencentcloud.com/search/?wd=DescribeSecurityGroupRules&type=document&lang=en)
* Update: [ModifySecurityGroupRule](https://www.tencentcloud.com/search/?wd=ModifySecurityGroupRule&type=document&lang=en)
* Delete: [DeleteSecurityGroupRule](https://www.tencentcloud.com/search/?wd=DeleteSecurityGroupRule&type=document&lang=en)

## Import

FWM security group rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_security_group_rule.example group-id-example#rule-id-example
```

