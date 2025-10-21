---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_resources"
sidebar_current: "docs-tencentcloud-datasource-clb_resources"
description: |-
  Use this data source to query detailed information of clb resources
---

# tencentcloud_clb_resources

Use this data source to query detailed information of clb resources

## Example Usage

```hcl
data "tencentcloud_clb_resources" "resources" {
  filters {
    name   = "isp"
    values = ["BGP"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter to query the list of AZ resources as detailed below: zone - String - Optional - Filter by AZ, such as ap-guangzhou-1. isp -- String - Optional - Filter by the ISP. Values: BGP, CMCC, CUCC and CTCC.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter value array.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zone_resource_set` - List of resources supported by the AZ.
  * `edge_zone` - Whether the AZ is an edge zone. Values: true, false.
  * `ip_version` - IP version. Values: IPv4, IPv6, and IPv6_Nat.
  * `local_zone` - Whether the AZ is a LocalZone. Values: true, false.
  * `master_zone` - Primary AZ, such as ap-guangzhou-1.
  * `resource_set` - List of resources. Note: This field may return null, indicating that no valid values can be obtained.
    * `availability_set` - Available resources. Note: This field may return null, indicating that no valid values can be obtaine.
      * `availability` - Whether the resource is available. Values: Available, Unavailable.
      * `type` - Specific ISP resource information. Values: CMCC, CUCC, CTCC, BGP.
    * `isp` - ISP information, such as CMCC, CUCC, CTCC, BGP, and INTERNAL.
    * `type` - Specific ISP resource information, Vaules: CMCC, CUCC, CTCC, BGP, and INTERNAL.
  * `slave_zone` - Secondary AZ, such as ap-guangzhou-2. Note: This field may return null, indicating that no valid values can be obtained.
  * `zone_region` - Region of the AZ, such as ap-guangzhou.
  * `zone_resource_type` - Type of resources in the zone. Values: SHARED, EXCLUSIVE.


