---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_bandwidth_limit"
sidebar_current: "docs-tencentcloud-resource-ccn_bandwidth_limit"
description: |-
  Provides a resource to limit CCN bandwidth.
---

# tencentcloud_ccn_bandwidth_limit

Provides a resource to limit CCN bandwidth.

## Example Usage

### Set the upper limit of regional outbound bandwidth

```hcl
variable "other_region1" {
  default = "ap-shanghai"
}

resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_ccn_bandwidth_limit" "limit1" {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  bandwidth_limit = 500
}
```

### Set the upper limit between regions

```hcl
variable "other_region1" {
  default = "ap-shanghai"
}

variable "other_region2" {
  default = "ap-nanjing"
}

resource tencentcloud_ccn main {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
  qos                  = "AG"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  dst_region      = var.other_region2
  bandwidth_limit = 100
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String, ForceNew) ID of the CCN.
* `region` - (Required, String, ForceNew) Limitation of region.
* `bandwidth_limit` - (Optional, Int) Limitation of bandwidth. Default is `0`.
* `dst_region` - (Optional, String, ForceNew) Destination area restriction. If the `CCN` rate limit type is `OUTER_REGION_LIMIT`, this value does not need to be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeCcns](https://www.tencentcloud.com/document/product/215/34787) - This API (DescribeCcns) is used to query the CCN list.
* Create: [SetCcnRegionBandwidthLimits](https://www.tencentcloud.com/document/product/215/34779) - This API (SetCcnRegionBandwidthLimits) is used to set the outbound bandwidth cap for CCNs in each region. This API can only set the outbound bandwidth cap for regions in the network instances that have already been associated.
* Read: [DescribeCcns](https://www.tencentcloud.com/document/product/215/34787) - This API (DescribeCcns) is used to query the CCN list.
* Read: [GetCcnRegionBandwidthLimits](https://www.tencentcloud.com/document/product/215/37319) - This API is used to query the bandwidth limits of a CCN instance. Monthly-subscribed CCNs only support Inter-region Bandwidth Limits, while the pay-as-you-go CCNs support both the Inter-region Bandwidth Limits and Region Outbound Bandwidth
* Update: [DescribeCcns](https://www.tencentcloud.com/document/product/215/34787) - This API (DescribeCcns) is used to query the CCN list.
* Update: [SetCcnRegionBandwidthLimits](https://www.tencentcloud.com/document/product/215/34779) - This API (SetCcnRegionBandwidthLimits) is used to set the outbound bandwidth cap for CCNs in each region. This API can only set the outbound bandwidth cap for regions in the network instances that have already been associated.
* Delete: [SetCcnRegionBandwidthLimits](https://www.tencentcloud.com/document/product/215/34779) - This API (SetCcnRegionBandwidthLimits) is used to set the outbound bandwidth cap for CCNs in each region. This API can only set the outbound bandwidth cap for regions in the network instances that have already been associated.

