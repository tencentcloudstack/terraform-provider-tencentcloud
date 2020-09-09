---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_domain_error_pages"
sidebar_current: "docs-tencentcloud-datasource-gaap_domain_error_pages"
description: |-
  Use this data source to query custom GAAP HTTP domain error page info list.
---

# tencentcloud_gaap_domain_error_pages

Use this data source to query custom GAAP HTTP domain error page info list.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `domain` - (Required) HTTP domain to be queried.
* `listener_id` - (Required) ID of the layer7 listener to be queried.
* `ids` - (Optional) List of the error page info ID to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `error_page_info_list` - An information list of error page info detail. Each element contains the following attributes:
  * `body` - New response body.
  * `clear_headers` - Response headers to be removed.
  * `domain` - HTTP domain.
  * `error_codes` - Original error codes.
  * `id` - ID of the error page info.
  * `listener_id` - ID of the layer7 listener.
  * `new_error_codes` - New error code.
  * `set_headers` - Response headers to be set.


