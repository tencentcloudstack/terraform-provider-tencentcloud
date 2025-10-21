---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_key"
sidebar_current: "docs-tencentcloud-resource-api_gateway_api_key"
description: |-
  Use this resource to create API gateway access key.
---

# tencentcloud_api_gateway_api_key

Use this resource to create API gateway access key.

## Example Usage

### Automatically generate key for API gateway access key.

```hcl
resource "tencentcloud_api_gateway_api_key" "example_auto" {
  secret_name = "tf_example_auto"
  status      = "on"
}
```

### Manually generate a secret key for API gateway access key.

```hcl
resource "tencentcloud_api_gateway_api_key" "example_manual" {
  secret_name       = "tf_example_manual"
  status            = "on"
  access_key_type   = "manual"
  access_key_id     = "28e287e340507fa147b2c8284dab542f"
  access_key_secret = "0198a4b8c3105080f4acd9e507599eff"
}
```

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, String, ForceNew) Custom key name.
* `access_key_id` - (Optional, String) User defined key ID, required when access_key_type is manual. The length is 5-50 characters, consisting of letters, numbers, and English underscores.
* `access_key_secret` - (Optional, String) The user-defined key must be passed when the access_key_type is manual. The length is 10-50 characters, consisting of letters, numbers, and English underscores.
* `access_key_type` - (Optional, String) Key type, supports both auto and manual (custom keys), defaults to auto.
* `status` - (Optional, String) Key status. Valid values: `on`, `off`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.


## Import

API gateway access key can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_api_key.test AKIDMZwceezso9ps5p8jkro8a9fwe1e7nzF2k50B
```

