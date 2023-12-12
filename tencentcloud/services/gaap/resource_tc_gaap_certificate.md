Provides a resource to create a certificate of GAAP.

Example Usage

```hcl
resource "tencentcloud_gaap_certificate" "foo" {
  type    = "BASIC"
  content = "test:tx2KGdo3zJg/."
  name    = "test_certificate"
}
```

Import

GAAP certificate can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_certificate.foo cert-d5y6ei3b
```