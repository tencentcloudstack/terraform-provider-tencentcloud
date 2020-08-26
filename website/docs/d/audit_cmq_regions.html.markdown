---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_audit_cmq_regions"
sidebar_current: "docs-tencentcloud-datasource-audit_cmq_regions"
description: |-
  Use this data source to query the region list supported by the audit cmq.
---

# tencentcloud_audit_cmq_regions

Use this data source to query the region list supported by the audit cmq.

## Example Usage

```hcl
data "tencentcloud_audit_cmq_region" "cmq_region" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cmq_region_list` - List of available zones supported by cmq.
  * `cmq_region_name` - Cmq region description.
  * `cmq_region` - Cmq region.


