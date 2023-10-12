---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_bind_api_apps_status"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_bind_api_apps_status"
description: |-
  Use this data source to query detailed information of apiGateway bind_api_apps_status
---

# tencentcloud_api_gateway_bind_api_apps_status

Use this data source to query detailed information of apiGateway bind_api_apps_status

## Example Usage

```hcl
data "tencentcloud_api_gateway_bind_api_apps_status" "example" {
  service_id = "service-nxz6yync"
  api_ids    = ["api-0cvmf4x4", "api-jvqlzolk"]
  filters {
    name   = "ApiAppId"
    values = ["app-krljp4wn"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_ids` - (Required, Set: [`String`]) Array of API IDs.
* `service_id` - (Required, String) Service ID.
* `filters` - (Optional, List) Filter conditions. Supports ApiAppId, Environment, KeyWord (can match name or ID).
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Filter value of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - List of APIs bound by the application.
  * `api_app_api_set` - Application bound API information array.
    * `api_app_id` - Application ID.
    * `api_app_name` - Application Name.
    * `api_id` - API ID.
    * `api_name` - API name.
    * `api_region` - Apis region.
    * `authorized_time` - Authorization binding time, expressed in accordance with the ISO8601 standard and using UTC time. The format is: YYYY-MM-DDThh:mm:ssZ.
    * `environment_name` - Authorization binding environment.
    * `service_id` - Service ID.


