---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_nat_acl_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_nat_acl_rule"
description: |-
  Provides a resource to create a FWM NAT ACL rule
---

# tencentcloud_fwm_nat_acl_rule

Provides a resource to create a FWM NAT ACL rule

## Example Usage

```hcl
resource "tencentcloud_fwm_nat_acl_rule" "example" {
  group_id       = "fwmrg-in3q0qkq0n"
  direction      = 0
  source_content = "0.0.0.0/0"
  source_type    = "net"
  target_content = "0.0.0.0/0"
  target_type    = "net"
  protocol       = "ANY"
  rule_action    = "accept"
  scope          = "ALL"
  port           = "-1/-1"
  description    = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, Int, ForceNew) Rule direction: 0-outbound, 1-inbound. It is also part of the resource id, so changing it will recreate the rule.
* `group_id` - (Required, String, ForceNew) Rule group ID.
* `protocol` - (Required, String) Protocol: TCP/UDP/ICMP/ANY/HTTP/HTTPS/DNS/FTP...
* `rule_action` - (Required, String) Action: accept/drop/log (lowercase).
* `scope` - (Required, String) Effective range: ALL/region/nat-instance-id.
* `source_content` - (Required, String) Source address content.
* `source_type` - (Required, String, ForceNew) Source type: ip/net/url/domain/template/addressbook/instance/tag/location/vendor/dnsparse.
* `target_content` - (Required, String) Target address content.
* `target_type` - (Required, String, ForceNew) Target type (same enum as source_type).
* `belong_member_id` - (Optional, String) Member account ID (required when source/target type in instance/tag/template/addressbook).
* `description` - (Optional, String) Rule description.
* `param_template_id` - (Optional, String) Port template ID.
* `port` - (Optional, String) Port e.g. 80, 80,443, 80-90; ANY = -1/-1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID, it is part of the resource id `group_id#rule_id#direction`.
* `sequence` - Rule priority sequence returned by the server.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateNatAclRule](https://www.tencentcloud.com/search/?wd=CreateNatAclRule&type=document&lang=en)
* Read: [DescribeNatAclRules](https://www.tencentcloud.com/search/?wd=DescribeNatAclRules&type=document&lang=en)
* Update: [ModifyNatAclRule](https://www.tencentcloud.com/search/?wd=ModifyNatAclRule&type=document&lang=en)
* Delete: [DeleteNatAclRule](https://www.tencentcloud.com/search/?wd=DeleteNatAclRule&type=document&lang=en)

## Import

FWM NAT ACL rule can be imported using the groupId#ruleId#direction, e.g.

```
terraform import tencentcloud_fwm_nat_acl_rule.example fwmrg-in3q0qkq0n#7520b898-30ba-415d-ab27-ff4542071ea9#0
```

