---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_function_version"
sidebar_current: "docs-tencentcloud-resource-scf_function_version"
description: |-
  Provides a resource to create a scf function_version
---

# tencentcloud_scf_function_version

Provides a resource to create a scf function_version

## Example Usage

```hcl
resource "tencentcloud_scf_function_version" "function_version" {
  function_name = "keep-1676351130"
  namespace     = "default"
  description   = "for-terraform-test"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String, ForceNew) Name of the released function.
* `description` - (Optional, String, ForceNew) Function description.
* `namespace` - (Optional, String, ForceNew) Function namespace.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `function_version` - Version of the released function.


## Import

scf function_version can be imported using the id, e.g.

```
terraform import tencentcloud_scf_function_version.function_version functionName#namespace#functionVersion
```

