---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_domain_error_page_info"
sidebar_current: "docs-tencentcloud-resource-gaap_domain_error_page_info"
description: |-
  Provide a resource to custom error page info for a GAAP HTTP domain.
---

# tencentcloud_gaap_domain_error_page_info

Provide a resource to custom error page info for a GAAP HTTP domain.

## Example Usage

```hcl
resource tencentcloud_gaap_proxy "foo" {
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
  proxy_id = tencentcloud_gaap_proxy.foo.id
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_domain_error_page_info "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = tencentcloud_gaap_http_domain.foo.domain
  error_codes = [404, 503]
  body        = "bad request"
}
```

## Argument Reference

The following arguments are supported:

* `body` - (Required, ForceNew) New response body.
* `domain` - (Required, ForceNew) HTTP domain.
* `error_codes` - (Required, ForceNew) Original error codes.
* `listener_id` - (Required, ForceNew) ID of the layer7 listener.
* `clear_headers` - (Optional, ForceNew) Response headers to be removed.
* `new_error_code` - (Optional, ForceNew) New error code.
* `set_headers` - (Optional, ForceNew) Response headers to be set.


