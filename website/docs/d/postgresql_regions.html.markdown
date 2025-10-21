---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_regions"
sidebar_current: "docs-tencentcloud-datasource-postgresql_regions"
description: |-
  Use this data source to query detailed information of postgresql regions
---

# tencentcloud_postgresql_regions

Use this data source to query detailed information of postgresql regions

## Example Usage

```hcl
data "tencentcloud_postgresql_regions" "regions" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region_set` - Region information set.
  * `region_id` - Region number.
  * `region_name` - Region name.
  * `region_state` - Availability status. UNAVAILABLE: unavailable, AVAILABLE: available.
  * `region` - Region abbreviation.
  * `support_international` - Whether the resource can be purchased in this region. Valid values: `0` (no), `1` (yes).Note: this field may return `null`, indicating that no valid values can be obtained.


