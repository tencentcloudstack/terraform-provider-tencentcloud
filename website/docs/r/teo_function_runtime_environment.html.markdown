---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_runtime_environment"
sidebar_current: "docs-tencentcloud-resource-teo_function_runtime_environment"
description: |-
  Provides a resource to create a teo teo_function_runtime_environment
---

# tencentcloud_teo_function_runtime_environment

Provides a resource to create a teo teo_function_runtime_environment

## Example Usage

```hcl
resource "tencentcloud_teo_function_runtime_environment" "teo_function_runtime_environment" {
  function_id = "ef-txx7fnua"
  zone_id     = "zone-2qtuhspy7cr6"

  environment_variables {
    key   = "test-a"
    type  = "string"
    value = "AAA"
  }
  environment_variables {
    key   = "test-b"
    type  = "string"
    value = "BBB"
  }
}
```

## Argument Reference

The following arguments are supported:

* `environment_variables` - (Required, List) The environment variable list.
* `function_id` - (Required, String, ForceNew) ID of the Function.
* `zone_id` - (Required, String, ForceNew) ID of the site.

The `environment_variables` object supports the following:

* `key` - (Required, String) The name of the variable, which is limited to alphanumeric characters and the special characters `@`, `.`, `-`, and `_`. It can have a maximum of 64 bytes and should not be duplicated.
* `type` - (Required, String) The type of the variable can have the following values:  - `string`: Represents a string type.  - `json`: Represents a JSON object type.
* `value` - (Required, String) The value of the variable, which is limited to a maximum of 5000 bytes. The default value is empty.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo teo_function_runtime_environment can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment zone_id#function_id
```

