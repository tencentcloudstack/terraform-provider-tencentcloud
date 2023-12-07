---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_ip_strategy"
sidebar_current: "docs-tencentcloud-resource-api_gateway_ip_strategy"
description: |-
  Use this resource to create IP strategy of API gateway.
---

# tencentcloud_api_gateway_ip_strategy

Use this resource to create IP strategy of API gateway.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "niceservice"
  protocol     = "http&https"
  service_desc = "your nice service"
  net_type     = ["INNER", "OUTER"]
  ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_ip_strategy" "test" {
  service_id    = tencentcloud_api_gateway_service.service.id
  strategy_name = "tf_test"
  strategy_type = "BLACK"
  strategy_data = "9.9.9.9"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, String, ForceNew) The ID of the API gateway service.
* `strategy_data` - (Required, String) IP address data.
* `strategy_name` - (Required, String, ForceNew) User defined strategy name.
* `strategy_type` - (Required, String, ForceNew) Blacklist or whitelist.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `strategy_id` - IP policy ID.


## Import

IP strategy of API gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_ip_strategy.test service-ohxqslqe#IPStrategy-q1lk8ud2
```

