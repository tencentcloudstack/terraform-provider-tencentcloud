---
subcategory: "Audit"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_audit_cos_regions"
sidebar_current: "docs-tencentcloud-datasource-audit_cos_regions"
description: |-
  Use this data source to query the cos region list supported by the audit.
---

# tencentcloud_audit_cos_regions

Use this data source to query the cos region list supported by the audit.

## Example Usage

```hcl
data "tencentcloud_audit_cos_regions" "foo" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `audit_cos_region_list` - List of available regions supported by audit cos.
  * `cos_region_name` - Cos region chinese name.
  * `cos_region` - Cos region.


