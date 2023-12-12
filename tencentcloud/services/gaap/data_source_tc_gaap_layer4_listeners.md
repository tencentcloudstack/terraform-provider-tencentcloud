Use this data source to query gaap layer4 listeners.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource "tencentcloud_gaap_layer4_listener" "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener"
  port            = 80
  realserver_type = "IP"
  proxy_id        = tencentcloud_gaap_proxy.foo.id
  health_check    = true
  interval        = 5
  connect_timeout = 2

  realserver_bind_set {
    id   = tencentcloud_gaap_realserver.foo.id
    ip   = tencentcloud_gaap_realserver.foo.ip
    port = 80
  }
}

data "tencentcloud_gaap_layer4_listeners" "foo" {
  protocol    = "TCP"
  proxy_id    = tencentcloud_gaap_proxy.foo.id
  listener_id = tencentcloud_gaap_layer4_listener.foo.id
}
```