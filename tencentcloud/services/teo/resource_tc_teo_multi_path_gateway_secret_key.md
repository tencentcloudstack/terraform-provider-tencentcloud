Provides a resource to manage a TEO multi-path gateway secret key config.

Example Usage

```hcl
resource "tencentcloud_teo_multi_path_gateway_secret_key" "example" {
  zone_id    = "zone-3edjdliiw3he"
  secret_key = "dGVzdC1zZWNyZXQta2V5LWZvci1tdWx0aS1wYXRoLWdhdGV3YXk="
}
```

Import

TEO multi-path gateway secret key can be imported using the id, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway_secret_key.example zone-3edjdliiw3he
```
