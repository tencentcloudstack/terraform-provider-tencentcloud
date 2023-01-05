---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_plugin_attachment"
sidebar_current: "docs-tencentcloud-resource-api_gateway_plugin_attachment"
description: |-
  Provides a resource to create a apiGateway plugin_attachment
---

# tencentcloud_api_gateway_plugin_attachment

Provides a resource to create a apiGateway plugin_attachment

## Example Usage

```hcl
resource "tencentcloud_api_gateway_plugin_attachment" "plugin_attachment" {
  plugin_id        = "plugin-ny74siyz"
  service_id       = "service-n1mgl0sq"
  environment_name = "test"
  api_id           = "api-6tfrdysk"
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required, String, ForceNew) Id of API.
* `environment_name` - (Required, String, ForceNew) Name of Environment.
* `plugin_id` - (Required, String, ForceNew) Id of Plugin.
* `service_id` - (Required, String, ForceNew) Id of Service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

apiGateway plugin_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_api_gateway_plugin_attachment.plugin_attachment pluginId#serviceId#environmentName#apiId
```

