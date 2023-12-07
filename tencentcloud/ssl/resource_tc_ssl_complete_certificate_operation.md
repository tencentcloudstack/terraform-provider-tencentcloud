Provides a resource to create a ssl complete_certificate

Example Usage

```hcl
resource "tencentcloud_ssl_complete_certificate_operation" "complete_certificate" {
  certificate_id = "9Bfe1IBR"
}
```

Import

ssl complete_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_complete_certificate_operation.complete_certificate complete_certificate_id
```