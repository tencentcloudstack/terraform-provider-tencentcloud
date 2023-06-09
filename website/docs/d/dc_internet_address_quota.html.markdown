---
subcategory: "Direct Connect(DC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_internet_address_quota"
sidebar_current: "docs-tencentcloud-datasource-dc_internet_address_quota"
description: |-
  Use this data source to query detailed information of dc internet_address_quota
---

# tencentcloud_dc_internet_address_quota

Use this data source to query detailed information of dc internet_address_quota

## Example Usage

```hcl
data "tencentcloud_dc_internet_address_quota" "internet_address_quota" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ipv4_bgp_num` - Number of used BGP type IPv4 Internet addresses.
* `ipv4_bgp_quota` - BGP type IPv4 Internet address quota.
* `ipv4_other_num` - The number of non-BGP Internet addresses used.
* `ipv4_other_quota` - Non-BGP type IPv4 Internet address quota.
* `ipv6_prefix_len` - The minimum prefix length allowed on the IPv6 Internet public network.


