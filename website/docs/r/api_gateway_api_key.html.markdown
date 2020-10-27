---
subcategory: "API GateWay"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_key"
sidebar_current: "docs-tencentcloud-resource-api_gateway_api_key"
description: |-
  Use this resource to create API gateway access key.
---

# tencentcloud_api_gateway_api_key

Use this resource to create API gateway access key.

## Example Usage

```hcl
resource "tencentcloud_api_gateway_api_key" "test" {
  secret_name = "my_api_key"
  status      = "on"
}
```

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, ForceNew) Custom key name.
* `status` - (Optional) Key status. Valid values: `on`, `off`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `access_key_secret` - Created API key.
* `create_time` - Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.
* `modify_time` - Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.


## Import

API gateway access key can be imported using the id, e.g.

```
$ terraform import tencentcloud_api_gateway_api_key.test AKIDMZwceezso9ps5p8jkro8a9fwe1e7nzF2k50B
```

