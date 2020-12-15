---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_ip_strategies"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_ip_strategies"
description: |-
  Use this data source to query API gateway IP strategy.
---

# tencentcloud_api_gateway_ip_strategies

Use this data source to query API gateway IP strategy.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_service" "service" {
  service_name = "ck"
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

data "tencentcloud_api_gateway_ip_strategies" "id" {
  service_id = tencentcloud_api_gateway_ip_strategy.test.service_id
}

data "tencentcloud_api_gateway_ip_strategies" "name" {
  service_id    = tencentcloud_api_gateway_ip_strategy.test.service_id
  strategy_name = tencentcloud_api_gateway_ip_strategy.test.strategy_name
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required) The service ID to be queried.
* `result_output_file` - (Optional) Used to save results.
* `strategy_name` - (Optional) Name of IP policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of strategy.
  * `attach_list` - List of bound API details.
    * `api_business_type` - The type of oauth API. This field is valid when the `auth_type` is `OAUTH`, and the values are `NORMAL` (business API) and `OAUTH` (authorization API).
    * `api_desc` - API interface description.
    * `api_id` - The API ID.
    * `api_name` - API name.
    * `api_type` - API type. Valid values: `NORMAL`, `TSF`. `NORMAL` means common API, `TSF` means microservice API.
    * `auth_relation_api_id` - The unique ID of the associated authorization API, which takes effect when the authType is `OAUTH` and `ApiBusinessType` is normal. Identifies the unique ID of the oauth2.0 authorization API bound to the business API.
    * `auth_type` - API authentication type. Valid values: `SECRET`, `NONE`, `OAUTH`. `SECRET` means key pair authentication, `NONE` means no authentication.
    * `create_time` - Creation time in the format of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used.
    * `method` - API request method.
    * `modify_time` - Last modified time in the format of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used.
    * `oauth_config` - OAUTH configuration information. It takes effect when authType is `OAUTH`.
    * `path` - API path.
    * `protocol` - API protocol.
    * `relation_business_api_ids` - List of business API associated with authorized API.
    * `service_id` - The service ID.
    * `tags` - The label information associated with the API.
    * `uniq_vpc_id` - VPC unique ID.
    * `vpc_id` - VPC ID.
  * `bind_api_total_count` - The number of API bound to the strategy.
  * `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `ip_list` - The list of IP.
  * `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
  * `service_id` - The service ID.
  * `strategy_id` - The strategy ID.
  * `strategy_name` - Name of the strategy.
  * `strategy_type` - Type of the strategy.


