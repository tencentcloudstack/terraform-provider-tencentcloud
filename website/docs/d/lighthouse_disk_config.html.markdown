---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_disk_config"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_disk_config"
description: |-
  Use this data source to query detailed information of lighthouse disk_config
---

# tencentcloud_lighthouse_disk_config

Use this data source to query detailed information of lighthouse disk_config

## Example Usage

```hcl
data "tencentcloud_lighthouse_disk_config" "disk_config" {
  filters {
    name   = "zone"
    values = ["ap-guangzhou-3"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter list.zoneFilter by availability zone.Type: StringRequired: no.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Filter value of field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `disk_config_set` - List of cloud disk configurations.
  * `disk_sales_state` - Cloud disk sale status.
  * `disk_step_size` - Cloud disk increment.
  * `disk_type` - Cloud disk type.
  * `max_disk_size` - Maximum cloud disk size.
  * `min_disk_size` - Minimum cloud disk size.
  * `zone` - Availability zone.


