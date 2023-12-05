Provides a resource to create a layer7 listener of GAAP.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_layer7_listener" "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = tencentcloud_gaap_proxy.foo.id
}
```

Import

GAAP layer7 listener can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_layer7_listener.foo listener-11112222
```