---
subcategory: "Direct Connect(DC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_internet_address_statistics"
sidebar_current: "docs-tencentcloud-datasource-dc_internet_address_statistics"
description: |-
  Use this data source to query detailed information of dc internet_address_statistics
---

# tencentcloud_dc_internet_address_statistics

Use this data source to query detailed information of dc internet_address_statistics

## Example Usage

```hcl
data "tencentcloud_dc_internet_address_statistics" "internet_address_statistics" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `internet_address_statistics` - Statistical Information List of Internet Public Network Addresses.
  * `region` - region.
  * `subnet_num` - Number of Internet public network addresses.


