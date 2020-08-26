---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_audit_cos_regions"
sidebar_current: "docs-tencentcloud-datasource-audit_cos_regions"
description: |-
  Use this data source to query scaling configuration information.
---

# tencentcloud_audit_cos_regions

Use this data source to query scaling configuration information.

## Example Usage

```hcl
data "tencentcloud_audit_cos_region" "cos_region" {
  website_type = "zh"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used to save results.
* `website_type` - (Optional) Site type. zh means China region, en means international region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cos_region_list` - List of available zones supported by cos.
  * `cos_region_name` - cos region description.
  * `cos_region` - cos region.


