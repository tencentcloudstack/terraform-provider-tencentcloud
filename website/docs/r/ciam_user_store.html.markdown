---
subcategory: "Customer Identity and Access Management(CIAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ciam_user_store"
sidebar_current: "docs-tencentcloud-resource-ciam_user_store"
description: |-
  Provides a resource to create a ciam user store
---

# tencentcloud_ciam_user_store

Provides a resource to create a ciam user store

## Example Usage

```hcl
resource "tencentcloud_ciam_user_store" "user_store" {
  user_pool_name = "tf_user_store"
  user_pool_desc = "for terraform test 123"
  user_pool_logo = "https://ciam-prd-1302490086.cos.ap-guangzhou.myqcloud.com/temporary/92630252a2c5422d9663db5feafd619b.png"
}
```

## Argument Reference

The following arguments are supported:

* `user_pool_name` - (Required, String) User Store Name.
* `user_pool_desc` - (Optional, String) User Store Description.
* `user_pool_logo` - (Optional, String) User Store Logo.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateUserStore](https://cloud.tencent.com/document/api/1441/94460) - 创建用户目录
* Read: [ListUserStore](https://cloud.tencent.com/document/api/1441/94458) - 查询用户目录列表
* Update: [UpdateUserStore](https://cloud.tencent.com/document/api/1441/94457) - 更新用户目录
* Delete: [DeleteUserStore](https://cloud.tencent.com/document/api/1441/94459) - 删除用户目录

## Import

ciam user_store can be imported using the id, e.g.

```
terraform import tencentcloud_ciam_user_store.user_store userStoreId
```

