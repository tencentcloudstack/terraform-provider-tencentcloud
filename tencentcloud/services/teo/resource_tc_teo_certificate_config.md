Provides a resource to create a teo certificate

Example Usage

```hcl
resource "tencentcloud_teo_certificate_config" "certificate" {
  host    = "test.tencentcloud-terraform-provider.cn"
  mode    = "eofreecert"
  zone_id = "zone-2o1t24kgy362"
}
```

Configure SSL certificate

```hcl
resource "tencentcloud_teo_certificate_config" "certificate" {
  host    = "test.tencentcloud-terraform-provider.cn"
  mode    = "sslcert"
  zone_id = "zone-2o1t24kgy362"

  server_cert_info {
    cert_id     = "8xiUJIJd"
  }
}
```

Import

teo certificate can be imported using the id, e.g.

```
terraform import tencentcloud_teo_certificate_config.certificate zone_id#host#cert_id
```