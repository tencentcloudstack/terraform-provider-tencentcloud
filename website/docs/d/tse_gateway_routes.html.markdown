---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_gateway_routes"
sidebar_current: "docs-tencentcloud-datasource-tse_gateway_routes"
description: |-
  Use this data source to query detailed information of tse gateway_routes
---

# tencentcloud_tse_gateway_routes

Use this data source to query detailed information of tse gateway_routes

## Example Usage

```hcl
data "tencentcloud_tse_gateway_routes" "gateway_routes" {
  gateway_id   = "gateway-ddbb709b"
  service_name = "test"
  route_name   = "keep-routes"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) gateway ID.
* `result_output_file` - (Optional, String) Used to save results.
* `route_name` - (Optional, String) route name.
* `service_name` - (Optional, String) service name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - result.
  * `route_list` - route list.
    * `created_time` - created time.
    * `destination_ports` - destination port for Layer 4 matching.
    * `force_https` - whether to enable forced HTTPS, no longer use.
    * `headers` - the headers of route.
      * `key` - key of header.
      * `value` - value of header.
    * `hosts` - host list.
    * `https_redirect_status_code` - https redirection status code.
    * `id` - service ID.
    * `methods` - method list.
    * `name` - service name.
    * `paths` - path list.
    * `preserve_host` - whether to keep the host when forwarding to the backend.
    * `protocols` - protocol list.
    * `service_id` - service ID.
    * `service_name` - service name.
    * `strip_path` - whether to strip path when forwarding to the backend.
  * `total_count` - total count.


