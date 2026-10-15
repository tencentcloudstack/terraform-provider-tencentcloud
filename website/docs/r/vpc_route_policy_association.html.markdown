---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_route_policy_association"
sidebar_current: "docs-tencentcloud-resource-vpc_route_policy_association"
description: |-
  Provides a resource to create a VPC route policy association
---

# tencentcloud_vpc_route_policy_association

Provides a resource to create a VPC route policy association

## Example Usage

```hcl
resource "tencentcloud_vpc_route_policy_association" "example" {
  route_policy_id = "rrp-7dnu4yoi"
  route_table_id  = "rtb-389phpuq"
  priority        = 10
}
```

## Argument Reference

The following arguments are supported:

* `priority` - (Required, Int) Priority.
* `route_policy_id` - (Required, String, ForceNew) Specifies the unique ID of the route reception policy.
* `route_table_id` - (Required, String, ForceNew) Unique route table ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateRoutePolicyAssociations](https://cloud.tencent.com/document/api/215/124429) - 创建路由接收策略绑定
* Read: [DescribeRoutePolicies](https://cloud.tencent.com/document/api/215/126466) - 用户查询路由策略列表
* Update: [ResetRoutePolicyAssociations](https://cloud.tencent.com/document/api/215/124420) - 重置路由策略绑定
* Delete: [DeleteRoutePolicyAssociations](https://cloud.tencent.com/document/api/215/124426) - 删除路由接收策略绑定

## Import

VPC route policy association can be imported using the routePolicyId#routeTableId, e.g.

```
terraform import tencentcloud_vpc_route_policy_association.example rrp-7dnu4yoi#rtb-389phpuq
```

