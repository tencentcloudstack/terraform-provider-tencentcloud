Provides a resource to manage a TEO multi-path gateway secret key config.

Example Usage

```hcl
resource "tencentcloud_teo_multi_path_gateway_secret_key" "example" {
  zone_id    = "zone-359h725djt7h"
  secret_key = base64encode("123123123")
}
```

Import

TEO multi-path gateway secret key can be imported using the id, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway_secret_key.example zone-3edjdliiw3he
```
