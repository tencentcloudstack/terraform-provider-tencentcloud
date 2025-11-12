Provides a resource to create a Organization external saml identity provider

~> **NOTE:** When creating it for the first time, you must set one of `encoded_metadata_document` and `x509_certificate`, `another_x509_certificate` cannot be set alone.

Example Usage

```hcl
resource "tencentcloud_organization_external_saml_identity_provider" "example" {
  zone_id                   = "z-1os7c9znogct"
  encoded_metadata_document = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz......RGVzY3JpcHRvcj4="
  another_x509_certificate  = <<-EOF
-----BEGIN CERTIFICATE-----
MIIC8DCCAdigAwIBAgIQPCotiH/l8K1K6kBgL4mBfzANBgkqhkiG9w0BAQsFADA0MTIwMAYDVQQD
......
qs39KP9jOtSzEzc1YhiX
-----END CERTIFICATE-----
EOF
}
```
