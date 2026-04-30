Provides a resource to manage the confirmation of TEO multi-path gateway origin ACL updates.

Example Usage

Confirm origin ACL version

```hcl
resource "tencentcloud_teo_confirm_multi_path_gateway_origin_acl" "example" {
  zone_id            = "zone-3edjdliiw3he"
  gateway_id         = "gw-abc12345"
  origin_acl_version = 2
}
```

Import

TEO confirm multi-path gateway origin ACL can be imported using the id, e.g.

```
terraform import tencentcloud_teo_confirm_multi_path_gateway_origin_acl.example zone-3edjdliiw3he#gw-abc12345
```
