---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_parse_resource_id"
sidebar_current: "docs-tencentcloud-function-parse_resource_id"
description: |-
  Provides a provider-defined function that parses a TencentCloud composite
resource ID into its constituent fields. The function is pure (no cloud
API call) and is safe to use in any Terraform expression context that
allows function calls.
---

# tencentcloud_parse_resource_id

Provides a provider-defined function that parses a TencentCloud composite
resource ID into its constituent fields. The function is pure (no cloud
API call) and is safe to use in any Terraform expression context that
allows function calls.

## Example Usage

```hcl
locals {
  parts = provider::tencentcloud::parse_resource_id("ins-abcd1234#vpc-xyz0987")
}

output "instance_id" {
  value = local.parts.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required, String) The composite resource id, for example "ins-abc#u-xyz".
* `separator` - (Required, String) The separator to split id by. Typically a single character, but any non-empty string is accepted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Splits an arbitrary composite resource id by a user-specified separator and returns the segments as a list of strings. Performs pure in-memory string processing; calls no TencentCloud API.


