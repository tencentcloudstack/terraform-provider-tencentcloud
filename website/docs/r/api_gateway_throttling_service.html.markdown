---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_throttling_service"
sidebar_current: "docs-tencentcloud-resource-api_gateway_throttling_service"
description: |-
  Use this resource to create API gateway throttling server.
---

# tencentcloud_api_gateway_throttling_service

Use this resource to create API gateway throttling server.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_throttling_service" "service" {
  service_id        = tencentcloud_api_gateway_service.service.id
  strategy          = "400"
  environment_names = ["release"]
}
```

## Argument Reference

The following arguments are supported:

* `environment_names` - (Required) List of Environment names.
* `service_id` - (Required, ForceNew) Service ID for query.
* `strategy` - (Required) Server QPS value. The service throttling value. Enter a positive number to limit the server query rate per second `QPS`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `environments` - A list of Throttling policy.
  * `environment_name` - Environment name.
  * `status` - Release status.
  * `strategy` - Throttling value.
  * `url` - Access service environment URL.
  * `version_name` - Published version number.


