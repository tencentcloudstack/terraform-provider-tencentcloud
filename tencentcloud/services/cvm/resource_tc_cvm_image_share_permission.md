Provides a resource to create a CVM image share permission

Example Usage

```hcl
resource "tencentcloud_cvm_image_share_permission" "example" {
  image_id    = "img-0elsru2u"
  account_ids = ["103849387508"]
}
```

Import

CVM image share permission can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_image_share_permission.example img-0elsru2u
```
