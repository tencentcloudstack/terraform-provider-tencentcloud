---
subcategory: "CynosDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_zone_config"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_zone_config"
description: |-
  Use this data source to query which instance types of Redis are available in a specific region.
---

# tencentcloud_cynosdb_zone_config

Use this data source to query which instance types of Redis are available in a specific region.

## Example Usage

```hcl
data "tencentcloud_cynosdb_zone_config" "foo" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of zone. Each element contains the following attributes:
  * `cpu` - Instance CPU, unit: core.
  * `machine_type` - Machine type.
  * `max_io_bandwidth` - Max io bandwidth.
  * `max_storage_size` - The maximum available storage for the instance, unit GB.
  * `memory` - Instance memory, unit: GB.
  * `min_storage_size` - Minimum available storage of the instance, unit: GB.
  * `zone_stock_infos` - Regional inventory information.
    * `has_stock` - Has stock.
    * `zone` - Availability zone.


