---
subcategory: "FireWall(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_strategy"
sidebar_current: "docs-tencentcloud-resource-fwm_strategy"
description: |-
  Provides a resource to create a FWM strategy
---

# tencentcloud_fwm_strategy

Provides a resource to create a FWM strategy

## Example Usage

```hcl
resource "tencentcloud_fwm_strategy" "example" {
  product         = "cfm"
  receive_account = ["10000xxx"]

  pre_strategy {
    group_id = "rule-group-id-1"
    sequence = 1
  }

  post_strategy {
    group_id = "rule-group-id-2"
    sequence = 1
  }

  receive_group = ["group-xxx"]
}
```

## Argument Reference

The following arguments are supported:

* `product` - (Required, String, ForceNew) Product type.
* `post_strategy` - (Optional, List) Post strategy rule group set.
* `pre_strategy` - (Optional, List) Pre strategy rule group set.
* `receive_account` - (Optional, List: [`String`]) Receive account list.
* `receive_group` - (Optional, List: [`String`]) Receive account group list.

The `post_strategy` object supports the following:

* `group_id` - (Required, String) Rule group ID.
* `sequence` - (Optional, Int) Priority.

The `pre_strategy` object supports the following:

* `group_id` - (Required, String) Rule group ID.
* `sequence` - (Optional, Int) Priority.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `strategy_ids` - Strategy ID list.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateStrategy](https://www.tencentcloud.com/search/?wd=CreateStrategy&type=document&lang=en)
  - [DescribeStrategy](https://www.tencentcloud.com/search/?wd=DescribeStrategy&type=document&lang=en)
* Read: [DescribeStrategy](https://www.tencentcloud.com/search/?wd=DescribeStrategy&type=document&lang=en)
* Update: [DescribeStrategy](https://www.tencentcloud.com/search/?wd=DescribeStrategy&type=document&lang=en)
* Delete: [DeleteStrategy](https://www.tencentcloud.com/search/?wd=DeleteStrategy&type=document&lang=en)

## Import

FWM strategy can be imported using the joined StrategyIds (joined by `#`), e.g. the create response returns StrategyIds `["s1","s2"]`, then import with `s1#s2`, e.g.

```
terraform import tencentcloud_fwm_strategy.example s1#s2
```

