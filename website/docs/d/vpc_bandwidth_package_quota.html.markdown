---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_bandwidth_package_quota"
sidebar_current: "docs-tencentcloud-datasource-vpc_bandwidth_package_quota"
description: |-
  Use this data source to query detailed information of vpc bandwidth_package_quota
---

# tencentcloud_vpc_bandwidth_package_quota

Use this data source to query detailed information of vpc bandwidth_package_quota

## Example Usage

```hcl
data "tencentcloud_vpc_bandwidth_package_quota" "bandwidth_package_quota" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `quota_set` - Bandwidth Package Quota Details.
  * `quota_current` - current amount.
  * `quota_id` - Quota type.
  * `quota_limit` - quota amount.


