---
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
resource "tencentcloud_ccn" "main"{
	name ="ci-temp-test-ccn"
	description="ci-temp-test-ccn-des"
	qos ="AG"
}

data "tencentcloud_ccn_bandwidth_limits" "limit" {
	ccn_id ="${tencentcloud_ccn.main.id}"
}

resource "tencentcloud_ccn_bandwidth_limit" "limit1" {
	ccn_id ="${tencentcloud_ccn.main.id}"
	region ="${var.other_region1}"
	bandwidth_limit = 500
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, ForceNew) ID of the CCN to be queried.
* `result_output_file` - (Optional, ForceNew) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `limits` - The bandwidth limits of regions
  * `bandwidth_limit` - Limitation of bandwidth.
  * `region` - Limitation of region.


