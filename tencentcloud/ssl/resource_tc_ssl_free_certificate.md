Provide a resource to create a Free Certificate.

~> **NOTE:** Once certificat created, it cannot be removed within 1 hours.

Example Usage

Currently, `package_type` only support type 2. 2=TrustAsia TLS RSA CA.

```hcl
resource "tencentcloud_ssl_free_certificate" "example" {
  dv_auth_method    = "DNS_AUTO"
  domain            = "example.com"
  package_type      = "2"
  contact_email     = "test@example.com"
  contact_phone     = "18352458901"
  validity_period   = 12
  csr_encrypt_algo  = "RSA"
  csr_key_parameter = "2048"
  csr_key_password  = "csr_pwd"
  alias             = "example_free_cert"
}
```

Import

FreeCertificate instance can be imported, e.g.
```
$ terraform import tencentcloud_ssl_free_certificate.test free_certificate-id
```