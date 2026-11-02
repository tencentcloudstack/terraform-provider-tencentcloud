---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_strategy"
sidebar_current: "docs-tencentcloud-resource-fwm_strategy"
description: |-
  Provides a resource to create a FWM (Firewall Manager) strategy.
---

# tencentcloud_fwm_strategy

Provides a resource to create a FWM (Firewall Manager) strategy.

## Example Usage

### Create pre strategy

```hcl
resource "tencentcloud_fwm_strategy" "example" {
  product = "cfw_edge_acl"
  receive_account = [
    "100011949846",
  ]
  receive_group = [
    "acg-y5tocz3d",
    "acg-e6cvbls1"
  ]

  pre_strategy {
    group_id = "fwmrg-qww1luux7z"
  }
}
```

### Create post strategy

```hcl
resource "tencentcloud_fwm_strategy" "example" {
  product = "cfw_edge_acl"
  receive_account = [
    "100011949846",
    "100049315301"
  ]
  receive_group = [
    "acg-y5tocz3d",
  ]

  post_strategy {
    group_id = "fwmrg-qww1luux7z"
  }
}
```

## Argument Reference

The following arguments are supported:

* `product` - (Required, String, ForceNew) Product type, e.g. enterprise_sg / cfw_vpc_acl / cfw_nat_acl / cfw_edge_acl.
* `post_strategy` - (Optional, List) Post Rule Group.
* `pre_strategy` - (Optional, List) Pre Rule Group.
* `receive_account` - (Optional, Set: [`String`]) List of receive account UINs.
* `receive_group` - (Optional, Set: [`String`]) List of receive account group IDs.

The `post_strategy` object supports the following:

* `group_id` - (Required, String) Rule group ID.

The `pre_strategy` object supports the following:

* `group_id` - (Required, String) Rule group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `strategy_id` - Strategy ID.

The `post_strategy` object exports the following:

* `sequence` - Priority.

The `pre_strategy` object exports the following:

* `sequence` - Priority.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateStrategy](https://www.tencentcloud.com/search/?wd=CreateStrategy&type=document&lang=en)
  - [DescribeStrategy](https://www.tencentcloud.com/search/?wd=DescribeStrategy&type=document&lang=en)
* Read: [DescribeStrategy](https://www.tencentcloud.com/search/?wd=DescribeStrategy&type=document&lang=en)
* Update:
  - [DescribeStrategy](https://www.tencentcloud.com/search/?wd=DescribeStrategy&type=document&lang=en)
  - [ModifyStrategy](https://www.tencentcloud.com/search/?wd=ModifyStrategy&type=document&lang=en)
* Delete: [DeleteStrategy](https://www.tencentcloud.com/search/?wd=DeleteStrategy&type=document&lang=en)

## Import

FWM strategy can be imported using the StrategyId, e.g.

```
terraform import tencentcloud_fwm_strategy.example fwms-s1d49yaf
```

