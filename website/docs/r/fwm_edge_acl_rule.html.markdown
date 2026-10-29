---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_edge_acl_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_edge_acl_rule"
description: |-
  Provides a resource to create a Edge ACL rule of Cloud Firewall (CFW)
---

# tencentcloud_fwm_edge_acl_rule

Provides a resource to create a Edge ACL rule of Cloud Firewall (CFW)

## Example Usage

```hcl
resource "tencentcloud_fwm_edge_acl_rule" "example" {
  group_id       = "fwmrg-9z03vakjj8"
  direction      = 1
  source_content = "0.0.0.0/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "TCP"
  port           = "-1/-1"
  rule_action    = "accept"
  scope          = "all"
  description    = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, Int, ForceNew) Rule direction: 0-egress, 1-ingress. Changing this will recreate the rule.
* `group_id` - (Required, String, ForceNew) Edge ACL rule group ID. The rule group is the parent container and must already exist before creating the rule.
* `belong_member_id` - (Optional, String) Member account ID that the rule belongs to (used in multi-account scenarios). Required when source_type or target_type is instance or tag.
* `description` - (Optional, String) Rule description, length 0-256 characters.
* `param_template_id` - (Optional, String) Parameter template ID.
* `port` - (Optional, String) Port. Supports single port, port range and comma separated multi ports, such as 80, 1-65535, 80,443,8080.
* `protocol` - (Optional, String) Protocol type: TCP, UDP, ICMP, ANY.
* `rule_action` - (Optional, String) Rule action: accept-pass, drop-block, log-observe.
* `scope` - (Optional, String) Effective scope: serial, side, all.
* `source_content` - (Optional, String) Source address content. Format depends on source_type: IP/CIDR for ip, domain for domain, template ID for template, instance ID list (comma separated) for instance, tag key:value for tag.
* `source_type` - (Optional, String, ForceNew) Source address type: ip, domain, template, instance, tag. Changing this will recreate the rule.
* `target_content` - (Optional, String) Target address content. Format is the same as source_content.
* `target_type` - (Optional, String, ForceNew) Target address type: ip, domain, template, instance, tag. Changing this will recreate the rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.
* `sequence` - Rule execution sequence returned by the server.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateEdgeAclRule](https://www.tencentcloud.com/search/?wd=CreateEdgeAclRule&type=document&lang=en)
* Read: [DescribeEdgeAclRules](https://www.tencentcloud.com/search/?wd=DescribeEdgeAclRules&type=document&lang=en)
* Update: [ModifyEdgeAclRule](https://www.tencentcloud.com/search/?wd=ModifyEdgeAclRule&type=document&lang=en)
* Delete: [DeleteEdgeAclRule](https://www.tencentcloud.com/search/?wd=DeleteEdgeAclRule&type=document&lang=en)

## Import

FWM edge acl rule can be imported using the groupId#ruleId#direction, e.g.

```
terraform import tencentcloud_fwm_edge_acl_rule.example fwmrg-9z03vakjj8#2fe0b492-5623-4646-b40e-7b0908307ca8#1
```

