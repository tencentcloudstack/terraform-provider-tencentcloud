---
subcategory: "Private Link(PLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_end_point_service_white_list"
sidebar_current: "docs-tencentcloud-resource-vpc_end_point_service_white_list"
description: |-
  Provides a resource to create a vpc end_point_service_white_list
---

# tencentcloud_vpc_end_point_service_white_list

Provides a resource to create a vpc end_point_service_white_list

## Example Usage

```hcl
resource "tencentcloud_vpc_end_point_service_white_list" "end_point_service_white_list" {
  user_uin             = "100020512675"
  end_point_service_id = "vpcsvc-69y13tdb"
  description          = "terraform for test"
}
```

## Argument Reference

The following arguments are supported:

* `end_point_service_id` - (Required, String) ID of endpoint service.
* `user_uin` - (Required, String) UIN.
* `description` - (Optional, String) Description of white list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create Time.
* `owner` - APPID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateVpcEndPointServiceWhiteList](https://cloud.tencent.com/document/api/215/54683) - 创建终端服务白名单
* Read: [DescribeVpcEndPointServiceWhiteList](https://cloud.tencent.com/document/api/215/54677) - 查询终端节点服务的服务白名单列表
* Update: [ModifyVpcEndPointServiceWhiteList](https://cloud.tencent.com/document/api/215/54672) - 修改终端节点服务白名单属性
* Delete: [DeleteVpcEndPointServiceWhiteList](https://cloud.tencent.com/document/api/215/54680) - 删除终端节点服务白名单

## Import

vpc end_point_service_white_list can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_end_point_service_white_list.end_point_service_white_list end_point_service_white_list_id
```

