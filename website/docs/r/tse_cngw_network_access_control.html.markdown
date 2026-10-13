---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_cngw_network_access_control"
sidebar_current: "docs-tencentcloud-resource-tse_cngw_network_access_control"
description: |-
  Provides a resource to create a tse cngw_network_access_control
---

# tencentcloud_tse_cngw_network_access_control

Provides a resource to create a tse cngw_network_access_control

## Example Usage

```hcl
resource "tencentcloud_tse_cngw_network_access_control" "cngw_network_access_control" {
  gateway_id = "gateway-cf8c99c3"
  group_id   = "group-a160d123"
  network_id = "network-372b1e84"
  access_control {
    mode            = "Whitelist"
    cidr_white_list = ["1.1.1.0"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) gateway ID.
* `group_id` - (Required, String, ForceNew) gateway group ID.
* `network_id` - (Required, String, ForceNew) network id.
* `access_control` - (Optional, List) access control policy.

The `access_control` object supports the following:

* `cidr_black_list` - (Optional, List) Black list.
* `cidr_white_list` - (Optional, List) White list.
* `mode` - (Optional, String) Access mode: `Whitelist`, `Blacklist`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribePublicNetwork](https://cloud.tencent.com/document/api/1364/103723) - 查询云原生API网关实例公网详情
* Create: [ModifyNetworkAccessStrategy](https://cloud.tencent.com/document/api/1364/102775) - 修改云原生API网关实例Kong访问策略
* Read: [DescribePublicNetwork](https://cloud.tencent.com/document/api/1364/103723) - 查询云原生API网关实例公网详情
* Update: [DescribePublicNetwork](https://cloud.tencent.com/document/api/1364/103723) - 查询云原生API网关实例公网详情
* Update: [ModifyNetworkAccessStrategy](https://cloud.tencent.com/document/api/1364/102775) - 修改云原生API网关实例Kong访问策略

## Import

tse cngw_route_rate_limit can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_network_access_control.cngw_network_access_control gatewayId#groupId#networkId
```

