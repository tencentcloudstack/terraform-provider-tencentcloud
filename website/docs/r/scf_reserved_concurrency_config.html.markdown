---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_reserved_concurrency_config"
sidebar_current: "docs-tencentcloud-resource-scf_reserved_concurrency_config"
description: |-
  Provides a resource to create a scf reserved_concurrency_config
---

# tencentcloud_scf_reserved_concurrency_config

Provides a resource to create a scf reserved_concurrency_config

## Example Usage

```hcl
resource "tencentcloud_scf_reserved_concurrency_config" "reserved_concurrency_config" {
  function_name            = "keep-1676351130"
  reserved_concurrency_mem = 128000
  namespace                = "default"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String, ForceNew) Specifies the function of which you want to configure the reserved quota.
* `reserved_concurrency_mem` - (Required, Int, ForceNew) Reserved memory quota of the function. Note: the upper limit for the total reserved quota of the function is the user's total concurrency memory minus 12800.
* `namespace` - (Optional, String, ForceNew) Function namespace. Default value: default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

scf reserved_concurrency_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_reserved_concurrency_config.reserved_concurrency_config reserved_concurrency_config_id
```

