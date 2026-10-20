---
subcategory: "Cloud Firewall(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_edge_acl_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_edge_acl_rule"
description: |-
  Provides a resource to create a FWM edge acl rule
---

# tencentcloud_fwm_edge_acl_rule

Provides a resource to create a FWM edge acl rule

## Example Usage

```hcl
resource "tencentcloud_fwm_edge_acl_rule" "example" {
  group_id = "group-xxxx"
  rules {
    order_index    = 1
    direction      = 1
    source_content = "1.1.1.1/32"
    source_type    = "ip"
    target_content = "0.0.0.0/0"
    target_type    = "ip"
    port           = "80"
    protocol       = "TCP"
    rule_action    = "accept"
    description    = "example rule"
    scope          = "serial"
  }
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, ForceNew) Rule group ID.
* `rules` - (Required, List) Rule list.

The `rules` object supports the following:

* `belong_member_id` - (Optional, String) Member account ID to which the rule belongs.
* `description` - (Optional, String) Rule description, 0-256 characters.
* `direction` - (Optional, Int) Rule direction: 0-outbound, 1-inbound.
* `order_index` - (Optional, Int) Rule execution order. The smaller the number, the higher the priority. When creating a rule group, it must start from 1 and strictly increase.
* `param_template_id` - (Optional, String) Parameter template ID.
* `port` - (Optional, String) Port, supports single port, port range and comma-separated multi-port.
* `protocol` - (Optional, String) Protocol type: TCP, UDP, ICMP, ANY.
* `rule_action` - (Optional, String) Rule action: accept-pass, drop-block, log-observe.
* `scope` - (Optional, String) Effective scope: serial, side, all.
* `source_content` - (Optional, String) Source address content.
* `source_type` - (Optional, String) Source address type: ip, domain, template, instance, tag.
* `target_content` - (Optional, String) Target address content.
* `target_type` - (Optional, String) Target address type: ip, domain, template, instance, tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_ids` - Rule ID list returned by the cloud API.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateEdgeAclRule](https://www.tencentcloud.com/search/?wd=CreateEdgeAclRule&type=document&lang=en)
* Read: [DescribeEdgeAclRules](https://www.tencentcloud.com/search/?wd=DescribeEdgeAclRules&type=document&lang=en)
* Update:
  - [CreateEdgeAclRule](https://www.tencentcloud.com/search/?wd=CreateEdgeAclRule&type=document&lang=en)
  - [DeleteEdgeAclRule](https://www.tencentcloud.com/search/?wd=DeleteEdgeAclRule&type=document&lang=en)
* Delete: [DeleteEdgeAclRule](https://www.tencentcloud.com/search/?wd=DeleteEdgeAclRule&type=document&lang=en)

## Import

FWM edge acl rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_edge_acl_rule.example group-xxxx#rule-xxxx
```

