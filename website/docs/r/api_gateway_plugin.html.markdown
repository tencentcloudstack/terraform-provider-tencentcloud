---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_plugin"
sidebar_current: "docs-tencentcloud-resource-api_gateway_plugin"
description: |-
  Provides a resource to create a apiGateway plugin
---

# tencentcloud_api_gateway_plugin

Provides a resource to create a apiGateway plugin

## Example Usage

```hcl
resource "tencentcloud_api_gateway_plugin" "plugin" {
  plugin_name = "terraform-plugin-test"
  plugin_type = "IPControl"
  plugin_data = jsonencode({
    "type" : "white_list",
    "blocks" : "1.1.1.1\n2.2.2.2",
  })
  description = "terraform test"
}
```

## Argument Reference

The following arguments are supported:

* `plugin_data` - (Required, String) Statement to define plugin.
* `plugin_name` - (Required, String) Name of the user define plugin. It must start with a letter and end with letter or number, the rest can contain letters, numbers and dashes(-). The length range is from 2 to 50.
* `plugin_type` - (Required, String) Type of plugin. Now support IPControl, TrafficControl, Cors, CustomReq, CustomAuth, Routing, TrafficControlByParameter, CircuitBreaker, ProxyCache.
* `description` - (Optional, String) Description of plugin.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

apiGateway plugin can be imported using the id, e.g.

```
terraform import tencentcloud_api_gateway_plugin.plugin plugin_id
```

