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

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, ForceNew) ID of the CCN.
* `region` - (Required, ForceNew) Limitation of region.
* `bandwidth_limit` - (Optional) Limitation of bandwidth.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



