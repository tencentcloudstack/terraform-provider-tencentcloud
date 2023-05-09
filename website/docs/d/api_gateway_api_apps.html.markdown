---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_apps"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_api_apps"
description: |-
  Use this data source to query list information of api_gateway api_app
---

# tencentcloud_api_gateway_api_apps

Use this data source to query list information of api_gateway api_app

## Example Usage

```hcl
data "tencentcloud_api_gateway_api_apps" "test" {
  api_app_id   = ["app-rj8t6zx3"]
  api_app_name = ["app_test"]
}
```

## Argument Reference

The following arguments are supported:

* `api_app_id` - (Optional, String) Api app ID.
* `api_app_name` - (Optional, String) Api app name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `api_app_list` - List of ApiApp.
  * `api_app_desc` - ApiApp description.
  * `api_app_id` - ApiApp ID.
  * `api_app_key` - ApiApp key.
  * `api_app_name` - ApiApp Name.
  * `api_app_secret` - ApiApp secret.
  * `created_time` - ApiApp create time.
  * `modified_time` - ApiApp modified time.


