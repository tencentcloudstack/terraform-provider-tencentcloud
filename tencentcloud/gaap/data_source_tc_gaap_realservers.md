Use this data source to query gaap realservers.

Example Usage

```hcl
resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

data "tencentcloud_gaap_realservers" "foo" {
  ip = tencentcloud_gaap_realserver.foo.ip
}
```