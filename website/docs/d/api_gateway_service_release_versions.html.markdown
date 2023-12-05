---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_service_release_versions"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_service_release_versions"
description: |-
  Use this data source to query detailed information of apiGateway service_release_versions
---

# tencentcloud_api_gateway_service_release_versions

Use this data source to query detailed information of apiGateway service_release_versions

## Example Usage

```hcl
data "tencentcloud_api_gateway_service_release_versions" "example" {
  service_id = "service-nxz6yync"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, String) The unique ID of the service to be queried.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - List of service releases.Note: This field may return null, indicating that no valid value can be obtained.
  * `version_desc` - Version description.Note: This field may return null, indicating that no valid value can be obtained.
  * `version_name` - Version number.Note: This field may return null, indicating that no valid value can be obtained.


