Use this data source to query GAAP certificate.

Example Usage

```hcl
resource "tencentcloud_gaap_certificate" "foo" {
  type    = "BASIC"
  content = "test:tx2KGdo3zJg/."
  name    = "test_certificate"
}

data "tencentcloud_gaap_certificates" "foo" {
  id = tencentcloud_gaap_certificate.foo.id
}
```