---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_sync_invoke_function"
sidebar_current: "docs-tencentcloud-resource-scf_sync_invoke_function"
description: |-
  Provides a resource to create a scf sync_invoke_function
---

# tencentcloud_scf_sync_invoke_function

Provides a resource to create a scf sync_invoke_function

## Example Usage

```hcl
resource "tencentcloud_scf_sync_invoke_function" "invoke_function" {
  function_name = "keep-1676351130"
  qualifier     = "2"
  namespace     = "default"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String, ForceNew) Function name.
* `event` - (Optional, String, ForceNew) Function running parameter, which is in the JSON format. Maximum parameter size is 6 MB. This field corresponds to event input parameter.
* `log_type` - (Optional, String, ForceNew) Valid value: None (default) or Tail. If the value is Tail, log in the response will contain the corresponding function execution log (up to 4KB).
* `namespace` - (Optional, String, ForceNew) Namespace. default is used if it's left empty.
* `qualifier` - (Optional, String, ForceNew) Version or alias of the function. It defaults to $DEFAULT.
* `routing_key` - (Optional, String, ForceNew) Traffic routing config in json format, e.g., {k:v}. Please note that both k and v must be strings. Up to 1024 bytes allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



