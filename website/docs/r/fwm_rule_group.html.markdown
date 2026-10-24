---
subcategory: "Tencent Cloud Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_rule_group"
sidebar_current: "docs-tencentcloud-resource-fwm_rule_group"
description: |-
  Provides a resource to create a FWM rule group
---

# tencentcloud_fwm_rule_group

Provides a resource to create a FWM rule group

## Example Usage

```hcl
resource "tencentcloud_fwm_rule_group" "example" {
  group_name = "tf-example"
  product    = "enterprise_sg"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required, String) Rule group name, cannot be empty, up to 50 characters.
* `product` - (Required, String, ForceNew) Product the rule group belongs to. Valid values: `enterprise_sg` (enterprise security group), `cfw_vpc_acl` (VPC border firewall), `cfw_nat_acl` (NAT border firewall), `cfw_edge_acl` (internet border firewall).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `app_id` - Tenant APP ID.
* `create_by` - Creator sub account UIN.
* `create_name` - Creator name.
* `create_time` - Creation time.
* `direction` - Direction. `1`: inbound, `0`: outbound, `-1`: no direction.
* `exec_areas` - Execution stages.
* `group_id` - Rule group ID.
* `receive_account` - Associated delivery accounts.
  * `members` - Members of the account group, only returned when `receiver_type` is `1`.
    * `app_id` - Member APP ID.
    * `member_id` - Member ID.
    * `nickname` - Member name.
    * `node_name` - Member node name.
    * `uin` - Member UIN.
  * `nickname` - Tenant name.
  * `receiver_type` - Receiver type. `0`: account UIN, `1`: account group.
  * `uin` - Tenant UIN.
* `rule_count` - Number of rules in the rule group.
* `update_by` - Updater sub account UIN.
* `update_name` - Updater name.
* `update_time` - Update time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateRuleGroup](https://www.tencentcloud.com/search/?wd=CreateRuleGroup&type=document&lang=en)
  - [DescribeRuleGroups](https://www.tencentcloud.com/search/?wd=DescribeRuleGroups&type=document&lang=en)
* Read: [DescribeRuleGroups](https://www.tencentcloud.com/search/?wd=DescribeRuleGroups&type=document&lang=en)
* Update: [ModifyRuleGroup](https://www.tencentcloud.com/search/?wd=ModifyRuleGroup&type=document&lang=en)
* Delete: [DeleteRuleGroup](https://www.tencentcloud.com/search/?wd=DeleteRuleGroup&type=document&lang=en)

## Import

FWM rule group can be imported using the id, e.g.

```
terraform import tencentcloud_fwm_rule_group.example fwmrg-wsrrhuvrib
```

