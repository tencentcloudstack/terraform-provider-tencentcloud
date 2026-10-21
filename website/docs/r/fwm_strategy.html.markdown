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

```hcl
resource "tencentcloud_fwm_strategy" "example" {
  product         = "enterprise_sg"
  sub_account_uin = "100028274628"
  exec_area       = "pre"
  group_id        = "fwmg-xxxxxx"
  sequence        = 1
  receive_account = ["100028274628"]
}
```

## Argument Reference

The following arguments are supported:

* `exec_area` - (Required, String, ForceNew) Execution area, pre or post. Determines whether the strategy is created in PreStrategy or PostStrategy.
* `product` - (Required, String, ForceNew) Product type, e.g. enterprise_sg / cfw_vpc_acl / cfw_nat_acl / cfw_edge_acl.
* `sub_account_uin` - (Required, String, ForceNew) Sub account UIN. Declarative parameter: the FWM cloud API does not expose a corresponding input field, so it is not sent on Create/Modify and is not read back. The value is retained in state.
* `group_id` - (Optional, String) Rule group ID.
* `receive_account` - (Optional, List: [`String`]) List of receive account UINs.
* `receive_group` - (Optional, List: [`String`]) List of receive account group IDs.
* `sequence` - (Optional, Int) Priority, starting from 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_by` - Creator.
* `create_name` - Creator name.
* `create_time` - Create time.
* `err_msg` - Dispatch failure reason.
* `error_type` - Dispatch failure reason type.
* `group_name` - Rule group name.
* `rule_count` - Rule count.
* `rule_status` - Strategy status.
* `strategy_id` - Strategy ID.
* `update_by` - Updater.
* `update_name` - Updater name.
* `update_time` - Update time.

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
terraform import tencentcloud_fwm_strategy.example fwms-xxxx
```

