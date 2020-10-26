---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_usage_plan"
sidebar_current: "docs-tencentcloud-resource-api_gateway_usage_plan"
description: |-
  Use this resource to create API gateway usage plan.
---

# tencentcloud_api_gateway_usage_plan

Use this resource to create API gateway usage plan.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_usage_plan" "plan" {
  usage_plan_name         = "my_plan"
  usage_plan_desc         = "nice plan"
  max_request_num         = 100
  max_request_num_pre_sec = 10
}
```

## Argument Reference

The following arguments are supported:

* `usage_plan_name` - (Required) Custom usage plan name.
* `max_request_num_pre_sec` - (Optional) Limit of requests per second. Valid values: -1, [1,2000]. The default value is -1, which indicates no limit.
* `max_request_num` - (Optional) Total number of requests allowed. Valid values: -1, [1,99999999]. The default value is -1, which indicates no limit.
* `usage_plan_desc` - (Optional) Custom usage plan description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `attach_api_keys` - Attach API keys list.
* `attach_list` - Attach service and API list.
  * `api_id` - The API ID, this value is empty if attach service.
  * `api_name` - The API name, this value is empty if attach service.
  * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `environment` - The environment name.
  * `method` - The API method, this value is empty if attach service.
  * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `path` - The API path, this value is empty if attach service.
  * `service_id` - The service ID.
  * `service_name` - The service name.
* `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.


## Import

API gateway usage plan can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_usage_plan.plan usagePlan-gyeafpab
```

