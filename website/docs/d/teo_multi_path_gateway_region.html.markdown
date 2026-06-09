---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_multi_path_gateway_region"
sidebar_current: "docs-tencentcloud-datasource-teo_multi_path_gateway_region"
description: |-
  Use this data source to query available regions of TEO multi-path gateway
---

# tencentcloud_teo_multi_path_gateway_region

Use this data source to query available regions of TEO multi-path gateway

## Example Usage

### Query multi-path gateway available regions by zone_id

```hcl
data "tencentcloud_teo_multi_path_gateway_region" "example" {
  zone_id = "zone-2noqxz9b6klw"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Site ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gateway_regions` - List of available gateway regions.
  * `cn_name` - Chinese name of the region.
  * `en_name` - English name of the region.
  * `region_id` - Region ID.


