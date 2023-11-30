Use this data source to query gaap proxies.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

data "tencentcloud_gaap_proxies" "foo" {
  ids = [tencentcloud_gaap_proxy.foo.id]
}
```