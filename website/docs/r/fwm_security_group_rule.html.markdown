---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_security_group_rule"
sidebar_current: "docs-tencentcloud-resource-fwm_security_group_rule"
description: |-
  Provides a resource to manage an enterprise security group rule of Tencent Cloud Firewall Manager (FWM).
---

# tencentcloud_fwm_security_group_rule

Provides a resource to manage an enterprise security group rule of Tencent Cloud Firewall Manager (FWM).

## Example Usage

```hcl
resource "tencentcloud_fwm_security_group_rule" "example" {
  group_name     = "tf-example-fwm-group"
  ip_version     = "4"
  source_content = "1.1.1.1/0"
  source_type    = "net"
  dest_content   = "0.0.0.0/0"
  dest_type      = "net"
  protocol       = "TCP"
  port           = "-1/-1"
  rule_action    = "accept"
  description    = "rule description."
}
```

### :

```hcl
resource "tencentcloud_fwm_security_group_rule" "example" {
  group_id       = "fwmrg-xxxxxxxx"
  ip_version     = "4"
  source_content = "1.1.1.1/0"
  source_type    = "net"
  dest_content   = "0.0.0.0/0"
  dest_type      = "net"
  protocol       = "TCP"
  port           = "-1/-1"
  rule_action    = "accept"
  description    = "rule description."
}
```

## Argument Reference

The following arguments are supported:

* `belong_member_id` - (Optional, String) ID of the member account the rule belongs to. Maps to `BelongMemberId` of `SecurityGroupRule`. The describe response returns a `MemberInfo` object (see the computed `belong_member`), so this attribute is not refreshed from the response.
* `description` - (Optional, String) Description of the rule. Maps to `Description` of `SecurityGroupRule`. The describe response returns the remark in `Detail` (see the computed `detail`), so this attribute is not refreshed from the response.
* `dest_content` - (Optional, String) Destination address of the rule. Maps to `DestContent` of `SecurityGroupRule`. The describe response returns it as `TargetId`, so it is not refreshed from the response.
* `dest_type` - (Optional, String) Destination type of the rule. Maps to `DestType` (string) of `SecurityGroupRule`. The describe response returns `TargetType` as an int (see the computed `target_type`), so this attribute is not refreshed from the response.
* `group_id` - (Optional, String, ForceNew) ID of the enterprise security group rule group. Maps to `GroupId` of `CreateSecurityGroupRule`, `DescribeSecurityGroupRules`, `ModifySecurityGroupRule` and `DeleteSecurityGroupRule`. Exactly one of `group_id` and `group_name` must be specified.
* `group_name` - (Optional, String, ForceNew) Name of a new enterprise security group rule group to create for this rule. Maps to `GroupName` of `CreateRuleGroup` (the product is fixed to `enterprise_sg`). The group is created on apply and deleted on destroy. Exactly one of `group_id` and `group_name` must be specified.
* `ip_version` - (Optional, String) IP version of the rule. Maps to `IpVersion` of `SecurityGroupRule`.
* `order_index` - (Optional, Int) Priority of the rule, the smaller the value the higher the priority. Maps to `OrderIndex` of `SecurityGroupRule`.
* `port` - (Optional, String) Port or port range of the rule. Maps to `Port` of `SecurityGroupRule`.
* `protocol_port_type` - (Optional, Int) Protocol port type of the rule. Maps to `ProtocolPortType` of `SecurityGroupRule`. The describe response uses a different encoding, so this attribute is not refreshed from the response.
* `protocol` - (Optional, String) Protocol of the rule, for example `TCP`, `UDP`, `ICMP` or `ANY`. Maps to `Protocol` of `SecurityGroupRule`.
* `rule_action` - (Optional, String) Action of the rule, for example `accept` or `drop`. Maps to `RuleAction` (string) of `SecurityGroupRule`. The describe response returns the action as an int `Strategy` (see the computed `strategy`), so this attribute is not refreshed from the response.
* `scope` - (Optional, String) Effective scope of the rule. Maps to `Scope` of `SecurityGroupRule`.
* `service_template_id` - (Optional, String) ID of the protocol port template used by the rule. Maps to `ServiceTemplateId` of `SecurityGroupRule`.
* `source_content` - (Optional, String) Source address of the rule. Maps to `SourceContent` of `SecurityGroupRule`. The API echoes it back as the int-encoded `SourceType`/`SourceId` pair, so it is not refreshed from the describe response.
* `source_type` - (Optional, String) Source type of the rule. Maps to `SourceType` (string) of `SecurityGroupRule`. The describe response returns `SourceType` as an int (see the computed `target_type`), so this attribute is not refreshed from the response.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `belong_member` - Member account the rule belongs to. Returned by the describe response as `BelongMember`.
  * `app_id` - AppId of the member. Maps to `AppId` of `MemberInfo`.
  * `member_id` - ID of the member, empty in personal scenarios. Maps to `MemberId` of `MemberInfo`.
  * `nickname` - Nickname of the member. Maps to `Nickname` of `MemberInfo`.
  * `node_name` - Department the member belongs to. Maps to `NodeName` of `MemberInfo`.
  * `uin` - Uin of the member. Maps to `Uin` of `MemberInfo`.
* `cidr` - Destination asset CIDR returned by the describe response as `Cidr`.
* `detail` - Rule remark returned by the describe response as `Detail`.
* `dns_parse_count` - Statistics of the IPs resolved from the domain name of the rule. Returned by the describe response as `DnsParseCount`.
  * `invalid_count` - Number of IPs not delivered. Maps to `InvalidCount` of `SgDnsParseCount`.
  * `valid_count` - Number of IPs successfully delivered. Maps to `ValidCount` of `SgDnsParseCount`.
* `instance_name` - Destination asset name returned by the describe response as `InstanceName`.
* `is_new` - Whether the rule was changed recently, `0` means no and `1` means yes. Returned by the describe response as `IsNew`.
* `parameter_name` - Destination template name returned by the describe response as `ParameterName`.
* `private_ip` - Destination asset private IP returned by the describe response as `PrivateIp`.
* `protocol_port_name` - Protocol port template name returned by the describe response as `ProtocolPortName`.
* `public_ip` - Destination asset public IP returned by the describe response as `PublicIp`.
* `region` - Region of the rule returned by the describe response as `Region`.
* `rule_id` - ID of the rule, returned by `CreateSecurityGroupRule` as `RuleIds[0]` and echoed by the describe response as `RuleId`.
* `rule_uuid` - Numeric rule ID (uuid), equal to `RuleId`. Returned by the describe response as `Id`.
* `sou_cidr` - Source asset CIDR returned by the describe response as `SouCidr`.
* `sou_instance_name` - Source asset name returned by the describe response as `SouInstanceName`.
* `sou_parameter_name` - Source template name returned by the describe response as `SouParameterName`.
* `sou_private_ip` - Source asset private IP returned by the describe response as `SouPrivateIp`.
* `sou_public_ip` - Source asset public IP returned by the describe response as `SouPublicIp`.
* `source_id` - Source rule content returned by the describe response as `SourceId`.
* `strategy` - Rule action returned by the describe response as `Strategy`.
* `target_id` - Destination rule content returned by the describe response as `TargetId`.
* `target_type` - Destination rule type returned by the describe response as `TargetType`.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateRuleGroup](https://www.tencentcloud.com/search/?wd=CreateRuleGroup&type=document&lang=en)
  - [CreateSecurityGroupRule](https://www.tencentcloud.com/search/?wd=CreateSecurityGroupRule&type=document&lang=en)
* Read: [DescribeSecurityGroupRules](https://www.tencentcloud.com/search/?wd=DescribeSecurityGroupRules&type=document&lang=en)
* Update: [ModifySecurityGroupRule](https://www.tencentcloud.com/search/?wd=ModifySecurityGroupRule&type=document&lang=en)
* Delete:
  - [DeleteRuleGroup](https://www.tencentcloud.com/search/?wd=DeleteRuleGroup&type=document&lang=en)
  - [DeleteSecurityGroupRule](https://www.tencentcloud.com/search/?wd=DeleteSecurityGroupRule&type=document&lang=en)

## Import

Tencent Cloud Firewall Manager enterprise security group rule can be imported using the composite id, which is `<group_id>#<rule_id>`, e.g.

```
terraform import tencentcloud_fwm_security_group_rule.example fwmrg-xxxxxxxx#fwmrule-yyyyyyyy
```

-> **Note:** Both segments of the composite id are required, because every read/update/delete call of this resource needs the rule group id together with the rule id.

