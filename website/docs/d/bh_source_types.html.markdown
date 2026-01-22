---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_source_types"
sidebar_current: "docs-tencentcloud-datasource-bh_source_types"
description: |-
  Use this data source to query detailed information of BH source types
---

# tencentcloud_bh_source_types

Use this data source to query detailed information of BH source types

## Example Usage

```hcl
data "tencentcloud_bh_source_types" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `source_type_set` - Authentication source information.
  * `name` - Account group source name.
  * `source` - Account group source.
  * `target` - Distinguish between ioa original and iam-mini.
  * `type` - Account group source type.


