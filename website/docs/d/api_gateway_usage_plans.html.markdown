---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_usage_plans"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_usage_plans"
description: |-
  Use this data source to query API gateway usage plans.
---

# tencentcloud_api_gateway_usage_plans

Use this data source to query API gateway usage plans.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_usage_plan" "plan" {
  usage_plan_name         = "my_plan"
  usage_plan_desc         = "nice plan"
  max_request_num         = 100
  max_request_num_pre_sec = 10
}

data "tencentcloud_api_gateway_usage_plans" "name" {
  usage_plan_name = tencentcloud_api_gateway_usage_plan.plan.usage_plan_name
}

data "tencentcloud_api_gateway_usage_plans" "id" {
  usage_plan_id = tencentcloud_api_gateway_usage_plan.plan.id
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `usage_plan_id` - (Optional, String) ID of the usage plan.
* `usage_plan_name` - (Optional, String) Name of the usage plan.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of usage plans.
  * `create_time` - Creation time in the format of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used.
  * `max_request_num_pre_sec` - Limit of requests per second. Valid values formats: `-1`, `[1,2000]`. The default value is -1, which indicates no limit.
  * `max_request_num` - Total number of requests allowed. Valid value formats: `-1`, `[1,99999999]`. The default value is -1, which indicates no limit.
  * `modify_time` - Last modified time in the format of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used.
  * `usage_plan_desc` - Custom usage plan description.
  * `usage_plan_id` - ID of the usage plan.
  * `usage_plan_name` - Name of the usage plan.


