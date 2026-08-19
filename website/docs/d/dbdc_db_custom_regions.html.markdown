---
subcategory: "Database Dedicated Cluster(DBDC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbdc_db_custom_regions"
sidebar_current: "docs-tencentcloud-datasource-dbdc_db_custom_regions"
description: |-
  Use this data source to query the supported DB Custom region list from the TencentCloud DBDC product.
---

# tencentcloud_dbdc_db_custom_regions

Use this data source to query the supported DB Custom region list from the TencentCloud DBDC product.

## Example Usage

### Query all dbdc db custom regions

```hcl
data "tencentcloud_dbdc_db_custom_regions" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region_set` - DB Custom supported region list.
  * `region_state` - Sale status. Values: SELL (normal sale), SOLD_OUT (sold out).
  * `region` - Region name.


