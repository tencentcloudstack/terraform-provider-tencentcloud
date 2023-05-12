---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_region"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_region"
description: |-
  Use this data source to query detailed information of lighthouse region
---

# tencentcloud_lighthouse_region

Use this data source to query detailed information of lighthouse region

## Example Usage

```hcl
data "tencentcloud_lighthouse_region" "region" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `region_set` - Region information list.
  * `is_china_mainland` - Whether the region is in the Chinese mainland.
  * `region_name` - Region description.
  * `region_state` - Region availability status.
  * `region` - Region name.


