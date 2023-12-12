Provides a resource to create a GAAP proxy.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"

  tags = {
    test = "test"
  }
}
```

Import

GAAP proxy can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_proxy.foo link-11112222
```