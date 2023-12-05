---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_keys"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_api_keys"
description: |-
  Use this data source to query API gateway access keys.
---

# tencentcloud_api_gateway_api_keys

Use this data source to query API gateway access keys.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_api_key" "test" {
  secret_name = "my_api_key"
  status      = "on"
}

data "tencentcloud_api_gateway_api_keys" "name" {
  secret_name = tencentcloud_api_gateway_api_key.test.secret_name
}

data "tencentcloud_api_gateway_api_keys" "id" {
  api_key_id = tencentcloud_api_gateway_api_key.test.id
}
```

## Argument Reference

The following arguments are supported:

* `api_key_id` - (Optional, String) Created API key ID, this field is exactly the same as ID.
* `result_output_file` - (Optional, String) Used to save results.
* `secret_name` - (Optional, String) Custom key name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of API keys.
  * `access_key_secret` - Created API key.
  * `api_key_id` - API key ID.
  * `create_time` - Creation time in the format of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used.
  * `modify_time` - Last modified time in the format of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used.
  * `status` - Key status. Values: `on`, `off`.


