---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_zones"
sidebar_current: "docs-tencentcloud-datasource-dbdc_db_custom_zones"
description: |-
  Use this data source to query available zones and their sale status from the TencentCloud DBDC product.
---

# tencentcloud_dbdc_db_custom_zones

Use this data source to query available zones and their sale status from the TencentCloud DBDC product.

## Example Usage

### Query all dbdc db custom zones

```hcl
data "tencentcloud_dbdc_db_custom_zones" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zone_set` - DB Custom available zone list.
  * `zone_state` - Zone sale status. Values: `SELL` (normal sale), `SOLD_OUT` (sold out).
  * `zone` - Available zone, such as `ap-guangzhou-3`.


