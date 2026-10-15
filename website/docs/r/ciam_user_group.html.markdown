---
subcategory: "Customer Identity and Access Management(CIAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ciam_user_group"
sidebar_current: "docs-tencentcloud-resource-ciam_user_group"
description: |-
  Provides a resource to create a ciam user group
---

# tencentcloud_ciam_user_group

Provides a resource to create a ciam user group

## Example Usage

```hcl
resource "tencentcloud_ciam_user_store" "user_store" {
  user_pool_name = "tf_user_store"
  user_pool_desc = "for terraform test"
  user_pool_logo = "https://ciam-prd-1302490086.cos.ap-guangzhou.myqcloud.com/temporary/92630252a2c5422d9663db5feafd619b.png"
}

resource "tencentcloud_ciam_user_group" "user_group" {
  display_name  = "tf_user_group"
  user_store_id = tencentcloud_ciam_user_store.user_store.id
  description   = "for terrafrom test"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required, String) User Group Name.
* `user_store_id` - (Required, String) User Store ID.
* `description` - (Optional, String) User Group Description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateUserGroup](https://cloud.tencent.com/document/api/1441/94455) - 创建用户组
* Read: [ListUserGroups](https://cloud.tencent.com/document/api/1441/94453) - 查询用户组列表
* Update: [UpdateUserGroup](https://cloud.tencent.com/document/api/1441/94452) - 更新用户组
* Delete: [DeleteUserGroups](https://cloud.tencent.com/document/api/1441/94454) - 批量删除用户组

## Import

ciam user_group can be imported using the id, e.g.

```
terraform import tencentcloud_ciam_user_group.user_group userStoreId#userGroupId
```

