---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_service_environment_list"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_service_environment_list"
description: |-
  Use this data source to query detailed information of apiGateway service_environment_list
---

# tencentcloud_api_gateway_service_environment_list

Use this data source to query detailed information of apiGateway service_environment_list

## Example Usage

```hcl
data "tencentcloud_api_gateway_service_environment_list" "example" {
  service_id = "service-nxz6yync"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, String) The unique ID of the service to be queried.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Service binding environment details.Note: This field may return null, indicating that no valid value can be obtained.
  * `environment_name` - Environment name.
  * `status` - Release status, 1 means released, 0 means not released.
  * `url` - Access path.
  * `version_name` - Running version.


