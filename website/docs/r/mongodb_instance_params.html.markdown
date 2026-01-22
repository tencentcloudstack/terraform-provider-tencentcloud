---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_params"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance_params"
description: |-
  Provides a resource to create a mongodb mongodb_instance_params
---

# tencentcloud_mongodb_instance_params

Provides a resource to create a mongodb mongodb_instance_params

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance_params" "mongodb_instance_params" {
  instance_id = "cmgo-xxxxxx"
  instance_params {
    key   = "cmgo.crossZoneLoadBalancing"
    value = "on"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `instance_params` - (Required, Set) Specify the parameter name and value to be modified.
* `modify_type` - (Optional, String) Operation types, including:
	- IMMEDIATELY: Adjust immediately;
	- DELAY: Delay adjustment;
Optional field. If this parameter is not configured, it defaults to immediate adjustment.

The `instance_params` object supports the following:

* `key` - (Required, String) Parameter names that need to be modified.
* `value` - (Required, String) The value corresponding to the parameter name to be modified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



