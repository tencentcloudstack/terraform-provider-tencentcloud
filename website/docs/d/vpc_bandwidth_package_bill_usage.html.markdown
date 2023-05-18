---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_bandwidth_package_bill_usage"
sidebar_current: "docs-tencentcloud-datasource-vpc_bandwidth_package_bill_usage"
description: |-
  Use this data source to query detailed information of vpc bandwidth_package_bill_usage
---

# tencentcloud_vpc_bandwidth_package_bill_usage

Use this data source to query detailed information of vpc bandwidth_package_bill_usage

## Example Usage

```hcl
data "tencentcloud_vpc_bandwidth_package_bill_usage" "bandwidth_package_bill_usage" {
  bandwidth_package_id = "bwp-234rfgt5"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_package_id` - (Required, String) The unique ID of the postpaid bandwidth package.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bandwidth_package_bill_bandwidth_set` - current billing amount.
  * `bandwidth_usage` - Current billing amount in Mbps.


