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

Set the upper limit of regional outbound bandwidth

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

Set the upper limit between regions

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

* `ccn_id` - (Required, ForceNew) ID of the CCN.
* `region` - (Required, ForceNew) Limitation of region.
* `bandwidth_limit` - (Optional) Limitation of bandwidth.
* `dst_region` - (Optional, ForceNew) Destination area restriction. If the `CCN` rate limit type is `OUTER_REGION_LIMIT`, this value does not need to be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



