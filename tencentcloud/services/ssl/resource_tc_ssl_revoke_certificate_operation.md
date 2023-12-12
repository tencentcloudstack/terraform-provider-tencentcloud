Provides a resource to create a ssl revoke_certificate

Example Usage

```hcl
resource "tencentcloud_ssl_revoke_certificate_operation" "revoke_certificate" {
  certificate_id = "7zUGkVab"
}
```

Import

ssl revoke_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_revoke_certificate_operation.revoke_certificate revoke_certificate_id
```