---
subcategory: "Cloud Firewall(FWM)"
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
  group_id       = "group-xxxx"
  edge_id        = "ALL"
  fw_group_id    = "ALL"
  ip_version     = "ipv4"
  source_type    = "ip"
  source_content = "10.0.0.0/24"
  dest_type      = "ip"
  dest_content   = "192.168.1.0/24"
  protocol       = "TCP"
  port           = "80"
  rule_action    = "accept"
  order_index    = 1
  description    = "test rule"
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
* `order_index` - (Required, Int) Priority (starting from 1).
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
* `all_total_count` - Total count without filtering.
* `app_id` - Member AppId.
* `create_time` - Rule creation time.
* `dest_name` - Destination asset name.
* `dns_parse_cnt` - Domain count.
* `fw_group_name` - Firewall group name.
* `hit_count` - Hit count.
* `invalid` - Whether the rule is invalid: 0 valid, 1 invalid.
* `last_hit_time` - Last hit time.
* `member_id` - Member ID.
* `nickname` - Member nickname.
* `node_name` - Department.
* `param_template_name` - Port template name.
* `rule_id` - Rule ID.
* `sequence` - Priority (read-only, same source as order_index).
* `source_name` - Source asset name.
* `uin` - Member Uin.
* `update_time` - Rule modification time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateVpcAclRule](https://www.tencentcloud.com/search/?wd=CreateVpcAclRule&type=document&lang=en)
* Read: [DescribeVpcAclRules](https://www.tencentcloud.com/search/?wd=DescribeVpcAclRules&type=document&lang=en)
* Update: [ModifyVpcAclRule](https://www.tencentcloud.com/search/?wd=ModifyVpcAclRule&type=document&lang=en)
* Delete: [DeleteVpcAclRule](https://www.tencentcloud.com/search/?wd=DeleteVpcAclRule&type=document&lang=en)

## Import

FWM VPC ACL rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_vpc_acl_rule.example group-xxxx#rule-xxxx
```

