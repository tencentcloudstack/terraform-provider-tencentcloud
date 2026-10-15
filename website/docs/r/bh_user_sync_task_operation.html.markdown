---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_user_sync_task_operation"
sidebar_current: "docs-tencentcloud-resource-bh_user_sync_task_operation"
description: |-
  Provides a resource to create a BH user sync task operation
---

# tencentcloud_bh_user_sync_task_operation

Provides a resource to create a BH user sync task operation

## Example Usage

```hcl
resource "tencentcloud_bh_user_sync_task_operation" "example" {
  user_kind = 1
}
```

## Argument Reference

The following arguments are supported:

* `user_kind` - (Required, Int, ForceNew) Synchronized user type, 1-synchronize IOA users.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSyncUserTask](https://cloud.tencent.com/document/api/1025/125474) - 创建用户同步任务
* Create: [DescribeUserSyncStatus](https://cloud.tencent.com/document/api/1025/125470) - 获取用户同步状态

