---
subcategory: "Firewall Manager(FWM)"
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
* `group_id` - Rule group ID.

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
terraform import tencentcloud_fwm_rule_group.example fwmrg-3r5c8tzb47
```

