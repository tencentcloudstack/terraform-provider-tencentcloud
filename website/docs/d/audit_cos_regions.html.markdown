---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_audit_cos_regions"
sidebar_current: "docs-tencentcloud-datasource-audit_cos_regions"
description: |-
  Use this data source to query the region list supported by the audit cos.
---

# tencentcloud_audit_cos_regions

Use this data source to query the region list supported by the audit cos.

## Example Usage

```hcl
data "tencentcloud_audit_cos_region" "cos_region" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cos_region_list` - List of available zones supported by cos.
  * `cos_region_name` - Cos region description.
  * `cos_region` - Cos region.


