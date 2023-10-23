---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_check_proxy_create"
sidebar_current: "docs-tencentcloud-datasource-gaap_check_proxy_create"
description: |-
  Use this data source to query detailed information of gaap check proxy create
---

# tencentcloud_gaap_check_proxy_create

Use this data source to query detailed information of gaap check proxy create

## Example Usage

```hcl
data "tencentcloud_gaap_check_proxy_create" "check_proxy_create" {
  access_region      = "Guangzhou"
  real_server_region = "Beijing"
  bandwidth          = 10
  concurrent         = 2
  ip_address_version = "IPv4"
  network_type       = "normal"
  package_type       = "Thunder"
  http3_supported    = 0
}
```

## Argument Reference

The following arguments are supported:

* `access_region` - (Required, String) The access (acceleration) area of the proxy. The value can be obtained through the interface DescribeAccessRegionsByDestRegion.
* `bandwidth` - (Required, Int) The upper limit of proxy bandwidth, in Mbps.
* `concurrent` - (Required, Int) The upper limit of chanproxynel concurrency, representing the number of simultaneous online connections, in tens of thousands.
* `real_server_region` - (Required, String) The origin area of the proxy. The value can be obtained through the interface DescribeDestRegions.
* `group_id` - (Optional, String) If creating a proxy under a proxy group, you need to fill in the ID of the proxy group.
* `ip_address_version` - (Optional, String) IP version, can be taken as IPv4 or IPv6, with a default value of IPv4.
* `network_type` - (Optional, String) Network type, can take values &amp;#39;normal&amp;#39;, &amp;#39;cn2&amp;#39;, default value normal.
* `package_type` - (Optional, String) Channel package type. Thunder represents the standard proxy group, Accelerator represents the game accelerator proxy, and CrossBorder represents the cross-border proxy.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `check_flag` - Query whether the proxy with the given configuration can be created, 1 can be created, 0 cannot be created.


