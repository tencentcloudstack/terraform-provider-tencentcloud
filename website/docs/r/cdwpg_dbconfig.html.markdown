---
subcategory: "CDWPG"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwpg_dbconfig"
sidebar_current: "docs-tencentcloud-resource-cdwpg_dbconfig"
description: |-
  Provides a resource to create a cdwpg cdwpg_dbconfig
---

# tencentcloud_cdwpg_dbconfig

Provides a resource to create a cdwpg cdwpg_dbconfig

## Example Usage

```hcl
resource "tencentcloud_cdwpg_dbconfig" "cdwpg_dbconfig" {
  instance_id = "cdwpg-ua8wkqrt"
  node_config_params {
    node_type       = "cn"
    parameter_name  = "log_min_duration_statement"
    parameter_value = "10001"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `node_config_params` - (Optional, Set) Node config parameters.

The `node_config_params` object supports the following:

* `node_type` - (Required, String) Node type.
* `parameter_name` - (Optional, String) Parameter name.
* `parameter_value` - (Optional, String) Parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



