---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_plugins"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_api_plugins"
description: |-
  Use this data source to query detailed information of apiGateway api_plugins
---

# tencentcloud_api_gateway_api_plugins

Use this data source to query detailed information of apiGateway api_plugins

## Example Usage

```hcl
data "tencentcloud_api_gateway_api_plugins" "example" {
  api_id           = "api-0cvmf4x4"
  service_id       = "service-nxz6yync"
  environment_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required, String) API ID to be queried.
* `service_id` - (Required, String) The service ID to be queried.
* `environment_name` - (Optional, String) Environment information.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - API list information that the plug-in can bind.
  * `attached_time` - Binding time.
  * `description` - Plugin description.
  * `environment` - Environment information.
  * `plugin_data` - Plug-in definition statement.
  * `plugin_id` - Plugin ID.
  * `plugin_name` - Plugin name.
  * `plugin_type` - Plugin type.


