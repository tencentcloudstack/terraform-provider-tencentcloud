Provides a resource to create a cvm image_share_permission

Example Usage

```hcl
resource "tencentcloud_cvm_image_share_permission" "image_share_permission" {
  image_id = "img-xxxxxx"
  account_ids = ["xxxxxx"]
}
```

Import

cvm image_share_permission can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_image_share_permission.image_share_permission image_share_permission_id
```