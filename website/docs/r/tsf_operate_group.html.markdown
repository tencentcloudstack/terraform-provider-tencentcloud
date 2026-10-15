---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_operate_group"
sidebar_current: "docs-tencentcloud-resource-tsf_operate_group"
description: |-
  Provides a resource to create a tsf operate_group
---

# tencentcloud_tsf_operate_group

Provides a resource to create a tsf operate_group

## Example Usage

```hcl
resource "tencentcloud_tsf_operate_group" "operate_group" {
  group_id = "group-ynd95rea"
  operate  = "start"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) group id.
* `operate` - (Required, String) Operation, `start`- start the group, `stop`- stop the group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeGroupAttribute](https://cloud.tencent.com/document/api/649/64979) - 获取部署组其他属性
* Create: [StartGroup](https://cloud.tencent.com/document/api/649/36057) - 启动虚拟机部署组
* Create: [StopGroup](https://cloud.tencent.com/document/api/649/36055) - 停止虚拟机部署组
* Read: [DescribeGroupAttribute](https://cloud.tencent.com/document/api/649/64979) - 获取部署组其他属性
* Update: [DescribeGroupAttribute](https://cloud.tencent.com/document/api/649/64979) - 获取部署组其他属性
* Update: [StartGroup](https://cloud.tencent.com/document/api/649/36057) - 启动虚拟机部署组
* Update: [StopGroup](https://cloud.tencent.com/document/api/649/36055) - 停止虚拟机部署组

