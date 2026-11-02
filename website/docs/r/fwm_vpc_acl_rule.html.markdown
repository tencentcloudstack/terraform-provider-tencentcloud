---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_vpc_acl_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_vpc_acl_rule"
description: |-
  Provides a resource to create a FWM VPC ACL rule
---

# tencentcloud_fwm_vpc_acl_rule

Provides a resource to create a FWM VPC ACL rule

## Example Usage

```hcl
resource "tencentcloud_fwm_vpc_acl_rule" "example" {
  group_id       = "fwmrg-y53zyyg7da"
  edge_id        = "all"
  fw_group_id    = "all"
  ip_version     = "ipv4"
  source_type    = "net"
  source_content = "0.0.0.0/0"
  dest_type      = "net"
  dest_content   = "0.0.0.0/0"
  protocol       = "ANY"
  port           = "-1/-1"
  rule_action    = "accept"
  description    = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required, String) Rule description, required and up to 30 characters.
* `dest_content` - (Required, String, ForceNew) Destination address content.
* `dest_type` - (Required, String, ForceNew) Destination type: ip/url/template/instance/tag.
* `edge_id` - (Required, String, ForceNew) Edge firewall ID: ALL or CFWS-xxx.
* `fw_group_id` - (Required, String, ForceNew) Firewall instance ID: ALL, cfwg-xxx or ccn-xxx.
* `group_id` - (Required, String, ForceNew) Rule group ID, the parent resource reference, required by all four interfaces, immutable.
* `ip_version` - (Required, String, ForceNew) IP version: ipv4 or ipv6. Not a Create input parameter, only returned by Describe.
* `protocol` - (Required, String) Protocol: TCP/UDP/ICMP/ANY etc.
* `rule_action` - (Required, String) Action: accept/drop/log (lowercase).
* `source_content` - (Required, String, ForceNew) Source address content.
* `source_type` - (Required, String, ForceNew) Source type: ip/url/template/instance/tag.
* `belong_member_id` - (Optional, String) Rule owner member account ID.
* `param_template_id` - (Optional, String) Port template ID.
* `port` - (Optional, String) Port, empty for ICMP protocol.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.
* `sequence` - Priority (read-only, same source as order_index).

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateVpcAclRule](https://www.tencentcloud.com/search/?wd=CreateVpcAclRule&type=document&lang=en)
* Read: [DescribeVpcAclRules](https://www.tencentcloud.com/search/?wd=DescribeVpcAclRules&type=document&lang=en)
* Update: [ModifyVpcAclRule](https://www.tencentcloud.com/search/?wd=ModifyVpcAclRule&type=document&lang=en)
* Delete: [DeleteVpcAclRule](https://www.tencentcloud.com/search/?wd=DeleteVpcAclRule&type=document&lang=en)

## Import

FWM VPC ACL rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_vpc_acl_rule.example fwmrg-y53zyyg7da#e42a7669-cf0c-42bb-b563-7cbf399a1f0f
```

