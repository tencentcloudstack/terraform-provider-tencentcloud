Provides a resource to create a Organization identity center external saml identity provider

Example Usage

```hcl
resource "tencentcloud_identity_center_external_saml_identity_provider" "example" {
  zone_id                   = "z-1os7c9znogct"
  sso_status                = "Enabled"
  encoded_metadata_document = "PD94bWwgdmVyc2lvbj0iM......VzY3JpcHRvcj4="
}
```

Or

```hcl
resource "tencentcloud_identity_center_external_saml_identity_provider" "example" {
  zone_id          = "z-1os7c9znogct"
  entity_id        = "https://sts.windows.net/d513d5bc-9f39-4069-ba9a-1eeab2ca58c1/"
  login_url        = "https://login.microsoftonline.com/d513d5bc-9f39-4069-ba9a-1eeab2ca58c1/saml2"
  sso_status       = "Enabled"
  x509_certificate = <<-EOF
-----BEGIN CERTIFICATE-----
MIIC8DCCAdigAwIBAgIQVbznAx6JSrhKG7gfJdx+jDANBgkqhkiG9w0BAQsFADA0MTIwMAYDVQQD
......
......
......
8hRskP2V6CH9PS0Zz2Zq
-----END CERTIFICATE-----
  EOF
}
```

Import

Organization identity center external saml identity provider can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_external_saml_identity_provider.example z-1os7c9znogct
```
