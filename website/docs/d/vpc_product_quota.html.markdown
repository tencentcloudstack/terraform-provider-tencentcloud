---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_product_quota"
sidebar_current: "docs-tencentcloud-datasource-vpc_product_quota"
description: |-
  Use this data source to query detailed information of vpc product_quota
---

# tencentcloud_vpc_product_quota

Use this data source to query detailed information of vpc product_quota

## Example Usage

```hcl
data "tencentcloud_vpc_product_quota" "product_quota" {
  product = "vpc"
}
```

## Argument Reference

The following arguments are supported:

* `product` - (Required, String) The name of the network product to be queried. The products that can be queried are:vpc, ccn, vpn, dc, dfw, clb, eip.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `product_quota_set` - ProductQuota Array.
  * `quota_current` - Current Quota.
  * `quota_id` - Quota Id.
  * `quota_limit` - Quota limit.
  * `quota_name` - Quota name.
  * `quota_region` - Quota region.


