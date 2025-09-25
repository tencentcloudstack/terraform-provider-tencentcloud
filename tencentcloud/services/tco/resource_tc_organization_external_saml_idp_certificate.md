Provides a resource to create a Organization external saml identity provider certificate

Example Usage

```hcl
resource "tencentcloud_organization_external_saml_idp_certificate" "example" {
  zone_id          = "z-dsj3ieme"
  x509_certificate = "MIIBtjCCAVugAwIBAgITBmyf1XSXNmY/Owua2eiedgPySjAKBggqhkj********"
}
```

Import

Organization external saml identity provider certificate can be imported using the id, e.g.

```
terraform import tencentcloud_organization_external_saml_idp_certificate.example z-dsj3ieme#idp-c-2jd8923je29dr34
```
