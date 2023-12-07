Provides a resource to create a ciam user store

Example Usage

```hcl
resource "tencentcloud_ciam_user_store" "user_store" {
  user_pool_name = "tf_user_store"
  user_pool_desc = "for terraform test 123"
  user_pool_logo = "https://ciam-prd-1302490086.cos.ap-guangzhou.myqcloud.com/temporary/92630252a2c5422d9663db5feafd619b.png"
}
```

Import

ciam user_store can be imported using the id, e.g.

```
terraform import tencentcloud_ciam_user_store.user_store userStoreId
```