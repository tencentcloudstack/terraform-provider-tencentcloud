---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_function_address"
sidebar_current: "docs-tencentcloud-datasource-scf_function_address"
description: |-
  Use this data source to query detailed information of scf function_address
---

# tencentcloud_scf_function_address

Use this data source to query detailed information of scf function_address

## Example Usage

```hcl
data "tencentcloud_scf_function_address" "function_address" {
  function_name = "keep-1676351130"
  namespace     = "default"
  qualifier     = "$LATEST"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String) Function name.
* `namespace` - (Optional, String) Function namespace.
* `qualifier` - (Optional, String) Function version.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `code_sha256` - SHA256 code of the function.
* `url` - Cos address of the function.


