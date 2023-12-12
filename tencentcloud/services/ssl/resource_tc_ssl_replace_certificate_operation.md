Provides a resource to create a ssl replace_certificate

Example Usage

```hcl
resource "tencentcloud_ssl_replace_certificate_operation" "replace_certificate" {
  certificate_id = "8L6JsWq2"
  valid_type = "DNS_AUTO"
  csr_type = "online"
}
```

Import

ssl replace_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_replace_certificate_operation.replace_certificate replace_certificate_id
```