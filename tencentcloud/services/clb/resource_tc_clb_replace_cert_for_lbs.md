Provides a resource to create a clb replace_cert_for_lbs

Example Usage

Replace Server Cert By Cert ID
```hcl
resource "tencentcloud_clb_replace_cert_for_lbs" "replace_cert_for_lbs" {
  old_certificate_id = "zjUMifFK"
  certificate {
    cert_id = "6vcK02GC"
  }
}
```

Replace Server Cert By Cert Content
```hcl
data "tencentcloud_ssl_certificates" "foo" {
  name = "keep-ssl-ca"
}

resource "tencentcloud_clb_replace_cert_for_lbs" "replace_cert_for_lbs" {
  old_certificate_id = data.tencentcloud_ssl_certificates.foo.certificates.0.id
  certificate {
    cert_name    = "tf-test-cert"
    cert_content = <<-EOT
-----BEGIN CERTIFICATE-----
xxxxxxxxxxxxxxxxxxxxxxxxxxx
-----END CERTIFICATE-----
EOT
    cert_key     = <<-EOT
-----BEGIN RSA PRIVATE KEY-----
xxxxxxxxxxxxxxxxxxxxxxxxxxxx
-----END RSA PRIVATE KEY-----
EOT
  }
}
```

Replace Client Cert By Cert Content
```hcl
resource "tencentcloud_clb_replace_cert_for_lbs" "replace_cert_for_lbs" {
  old_certificate_id = "zjUMifFK"
  certificate {
    cert_ca_name = "tf-test-cert"
    cert_ca_content = <<-EOT
-----BEGIN CERTIFICATE-----
xxxxxxxxContentxxxxxxxxxxxxxx
-----END CERTIFICATE-----
EOT
  }
}
```

```
terraform import tencentcloud_clb_replace_cert_for_lbs.replace_cert_for_lbs replace_cert_for_lbs_id
```