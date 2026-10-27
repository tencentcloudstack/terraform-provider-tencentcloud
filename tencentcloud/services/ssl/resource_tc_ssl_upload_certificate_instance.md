Provides a resource to create a SSL upload certificate instance. The `tencentcloud_ssl_upload_certificate_instance` resource uploads new certificate public/private key content and updates associated cloud resources (such as clb) without changing the certificate id.



Example Usage

```hcl
resource "tencentcloud_ssl_upload_certificate_instance" "example" {
  old_certificate_id       = "f8xnynbd"
  resource_types           = ["clb"]
  certificate_public_key   = "-----BEGIN CERTIFICATE-----\nxxx\n-----END CERTIFICATE-----"
  certificate_private_key = "-----BEGIN RSA PRIVATE KEY-----\nxxx\n-----END RSA PRIVATE KEY-----"
  resource_types_regions {
    resource_type = "clb"
    regions       = ["ap-guangzhou"]
  }
}
```

Upload with certificate id

```hcl
resource "tencentcloud_ssl_upload_certificate_instance" "example" {
  old_certificate_id = "f8xnynbd"
  certificate_id      = "f8xnynbe"
  resource_types      = ["clb"]
}
```
