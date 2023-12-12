Provides a resource to create a ciam user group

Example Usage

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

Import

ciam user_group can be imported using the id, e.g.

```
terraform import tencentcloud_ciam_user_group.user_group userStoreId#userGroupId
```