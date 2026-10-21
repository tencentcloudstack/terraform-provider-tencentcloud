---
subcategory: "Firewall Managed(FWM)"
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
  group_id       = "fwmrg-abc"
  direction      = 0
  source_content = "1.1.1.1"
  source_type    = "ip"
  target_content = "0.0.0.0/0"
  target_type    = "ip"
  protocol       = "TCP"
  rule_action    = "accept"
  order_index    = 1
  scope          = "ALL"
  port           = "80"
  description    = "tf example nat acl rule"
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, Int, ForceNew) Rule direction: 0-outbound, 1-inbound.
* `group_id` - (Required, String, ForceNew) Rule group ID.
* `order_index` - (Required, Int) Priority, >=1, per-direction numbering.
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
* `app_id` - Member AppId (from BelongMember).
* `city_code` - City code.
* `city_key` - City key.
* `city_name` - City name.
* `cloud_code` - Cloud vendor code.
* `country_code` - Country code.
* `country_key` - Country key.
* `country_name` - Country name.
* `create_time` - Create time.
* `dns_parse_cnt` - Domain count.
* `fw_ins_id` - Firewall instance ID.
* `hit_count` - Hit count.
* `invalid` - Whether rule is invalid: 0-valid, 1-invalid.
* `is_cloud` - Whether cloud vendor rule: 0-no, 1-yes.
* `is_region` - Whether region rule: 0-no, 1-yes.
* `last_hit_time` - Last hit time.
* `member_id` - Member ID (from BelongMember).
* `nickname` - Member nickname (from BelongMember).
* `node_name` - Member department (from BelongMember).
* `param_template_name` - Port template name (display only).
* `rule_id` - Rule ID, part of composite id.
* `scope_desc` - Effective range description (language-dependent display).
* `sequence` - Sequence, display mirror of order_index.
* `source_name` - Source asset name (display only).
* `target_name` - Target asset name (display only).
* `uin` - Member Uin (from BelongMember).
* `update_time` - Update time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateNatAclRule](https://www.tencentcloud.com/search/?wd=CreateNatAclRule&type=document&lang=en)
* Read: [DescribeNatAclRules](https://www.tencentcloud.com/search/?wd=DescribeNatAclRules&type=document&lang=en)
* Update: [ModifyNatAclRule](https://www.tencentcloud.com/search/?wd=ModifyNatAclRule&type=document&lang=en)
* Delete: [DeleteNatAclRule](https://www.tencentcloud.com/search/?wd=DeleteNatAclRule&type=document&lang=en)

## Import

FWM NAT ACL rule can be imported using the groupId#ruleId, e.g.

```
terraform import tencentcloud_fwm_nat_acl_rule.example fwmrg-abc#fwmrule-xyz
```

