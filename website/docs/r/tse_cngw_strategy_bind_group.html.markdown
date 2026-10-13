---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_cngw_strategy_bind_group"
sidebar_current: "docs-tencentcloud-resource-tse_cngw_strategy_bind_group"
description: |-
  Provides a resource to create a tse cngw_strategy_bind_group
---

# tencentcloud_tse_cngw_strategy_bind_group

Provides a resource to create a tse cngw_strategy_bind_group

## Example Usage

```hcl
resource "tencentcloud_tse_cngw_strategy_bind_group" "cngw_strategy_bind_group" {
  gateway_id  = "gateway-cf8c99c3"
  strategy_id = "strategy-806ea0dd"
  group_id    = "group-a160d123"
  option      = "bind"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) gateway ID.
* `group_id` - (Required, String, ForceNew) group ID.
* `option` - (Required, String) `bind` or `unbind`.
* `strategy_id` - (Required, String, ForceNew) strategy ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Binding status.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [BindAutoScalerResourceStrategyToGroups](https://cloud.tencent.com/document/api/1364/103728) - 弹性伸缩策略批量绑定网关分组
* Create: [UnbindAutoScalerResourceStrategyFromGroups](https://cloud.tencent.com/document/api/1364/103721) - 弹性伸缩策略批量解绑网关分组
* Read: [DescribeAutoScalerResourceStrategyBindingGroups](https://cloud.tencent.com/document/api/1364/103724) - 查看弹性伸缩策略绑定的网关分组
* Update: [BindAutoScalerResourceStrategyToGroups](https://cloud.tencent.com/document/api/1364/103728) - 弹性伸缩策略批量绑定网关分组
* Update: [UnbindAutoScalerResourceStrategyFromGroups](https://cloud.tencent.com/document/api/1364/103721) - 弹性伸缩策略批量解绑网关分组

## Import

tse cngw_strategy_bind_group can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_strategy_bind_group.cngw_strategy_bind_group cngw_strategy_bind_group_id
```

