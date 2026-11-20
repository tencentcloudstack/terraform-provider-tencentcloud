---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_dispense_region"
sidebar_current: "docs-tencentcloud-data_source-monitor_dispense_region"
description: |-
  Use this data source to query the dispense (forwarding) deployment regions of the Monitor product.
---

# tencentcloud_monitor_dispense_region

Use this data source to query the dispense (forwarding) deployment regions of the Monitor product.

## Example Usage

```hcl
data "tencentcloud_monitor_dispense_region" "example" {}
```

## Argument Reference

The following arguments are supported:



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `region_cn_name` - Chinese name of the region.
  * `region_en_name` - English name of the region.
  * `region` - Region abbreviation, e.g. `ap-guangzhou`.
  * `rule_number` - Number of rules deployed in this region.
* `id` - ID of the data source, a random token.
* `region_list` - List of dispense deployment regions.


