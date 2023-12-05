Provides a resource to create a ssl download_certificate

Example Usage

```hcl
resource "tencentcloud_ssl_download_certificate_operation" "download_certificate" {
  certificate_id = "8x1eUSSl"
  output_path = "./"
}
```

Import

ssl download_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_download_certificate_operation.download_certificate download_certificate_id
```