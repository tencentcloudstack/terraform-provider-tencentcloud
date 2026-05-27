---
subcategory: "Global Accelerator 2(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_accelerate_regions"
sidebar_current: "docs-tencentcloud-datasource-ga2_accelerate_regions"
description: |-
  Use this data source to query available accelerate regions of GA2 (Global Accelerator 2)
---

# tencentcloud_ga2_accelerate_regions

Use this data source to query available accelerate regions of GA2 (Global Accelerator 2)

## Example Usage

### Query all GA2 accelerate regions

```hcl
data "tencentcloud_ga2_accelerate_regions" "example" {}
```

### Query GA2 accelerate regions and output to file

```hcl
data "tencentcloud_ga2_accelerate_regions" "example" {
  result_output_file = "ga2_accelerate_regions.json"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `accelerator_region_set` - Accelerate region list.
  * `area_name` - Area name.
  * `is_available` - Availability status. 0: unavailable, 1: available.
  * `is_china_mainland` - Whether it is a China mainland region.
  * `is_tencent_region` - Whether it is a Tencent region.
  * `name` - Region Chinese name.
  * `region` - Region identifier.
  * `support_isp_type` - Supported ISP types.


