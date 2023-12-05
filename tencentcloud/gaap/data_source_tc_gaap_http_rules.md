Use this data source to query forward rule of layer7 listeners.

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

resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource "tencentcloud_gaap_http_rule" "foo" {
  listener_id     = tencentcloud_gaap_layer7_listener.foo.id
  domain          = "www.qq.com"
  path            = "/"
  realserver_type = "IP"
  health_check    = true

  realservers {
    id   = tencentcloud_gaap_realserver.foo.id
    ip   = tencentcloud_gaap_realserver.foo.ip
    port = 80
  }
}

data "tencentcloud_gaap_http_rules" "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = tencentcloud_gaap_http_rule.foo.domain
}
```