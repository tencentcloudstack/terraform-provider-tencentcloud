---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_bandwidth_limits"
sidebar_current: "docs-tencentcloud-datasource-ccn_bandwidth_limits"
description: |-
  Use this data source to query detailed information of CCN bandwidth limits.
---

# tencentcloud_ccn_bandwidth_limits

Use this data source to query detailed information of CCN bandwidth limits.

## Example Usage

```hcl
variable "other_region1" {
  default = "ap-shanghai"
}

resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

data "tencentcloud_ccn_bandwidth_limits" "limit" {
  ccn_id = tencentcloud_ccn.main.id
}

resource "tencentcloud_ccn_bandwidth_limit" "limit1" {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  bandwidth_limit = 500
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String) ID of the CCN to be queried.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `limits` - The bandwidth limits of regions:
  * `bandwidth_limit` - Limitation of bandwidth.
  * `dst_region` - Destination area restriction.
  * `region` - Limitation of region.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [GetCcnRegionBandwidthLimits](https://cloud.tencent.com/document/api/215/45482) - This API is used to query the bandwidth limits of a CCN instance. Monthly-subscribed CCNs only support Inter-region Bandwidth Limits, while the pay-as-you-go CCNs support both the Inter-region Bandwidth Limits and Region Outbound Bandwidth ...

