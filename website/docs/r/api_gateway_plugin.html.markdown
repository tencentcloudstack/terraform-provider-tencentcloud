---
subcategory: "API GateWay(apigateway)"
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
resource "tencentcloud_api_gateway_plugin" "example" {
  plugin_name = "tf-example"
  plugin_type = "IPControl"
  plugin_data = jsonencode({
    "type" : "white_list",
    "blocks" : "1.1.1.1",
  })
  description = "desc."
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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreatePlugin](https://www.tencentcloud.com/search/?wd=CreatePlugin&type=document&lang=en)
* Read: [DescribePlugins](https://www.tencentcloud.com/search/?wd=DescribePlugins&type=document&lang=en)
* Update: [ModifyPlugin](https://www.tencentcloud.com/search/?wd=ModifyPlugin&type=document&lang=en)
* Delete: [DeletePlugin](https://www.tencentcloud.com/search/?wd=DeletePlugin&type=document&lang=en)

## Import

apiGateway plugin can be imported using the id, e.g.

```
terraform import tencentcloud_api_gateway_plugin.plugin plugin_id
```

