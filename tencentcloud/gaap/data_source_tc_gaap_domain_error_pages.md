Use this data source to query custom GAAP HTTP domain error page info list.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "%s"
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_domain_error_page "foo" {
  listener_id    = tencentcloud_gaap_layer7_listener.foo.id
  domain         = tencentcloud_gaap_http_domain.foo.domain
  error_codes    = [406, 504]
  new_error_code = 502
  body           = "bad request"
  clear_headers  = ["Content-Length", "X-TEST"]

  set_headers = {
    "X-TEST" = "test"
  }
}

data tencentcloud_gaap_domain_error_pages "foo" {
  listener_id = tencentcloud_gaap_domain_error_page.foo.listener_id
  domain      = tencentcloud_gaap_domain_error_page.foo.domain
}
```