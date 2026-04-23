---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_ip_region"
sidebar_current: "docs-tencentcloud-datasource-teo_ip_region"
description: |-
  Use this data source to query detailed information of TEO IP region
---

# tencentcloud_teo_ip_region

Use this data source to query detailed information of TEO IP region

## Example Usage

### Query IP region info

```hcl
data "tencentcloud_teo_ip_region" "example" {
  ips = [
    "1.1.1.1",
    "2.2.2.2"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `ips` - (Required, List: [`String`]) List of IP addresses (IPv4/IPv6) to query, up to 100 entries.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ip_region_info` - IP region information list.
  * `ip` - IP address, IPv4 or IPv6.
  * `is_edge_one_ip` - Whether the IP belongs to an EdgeOne node. Values: `yes` (belongs to EdgeOne node), `no` (does not belong to EdgeOne node).


