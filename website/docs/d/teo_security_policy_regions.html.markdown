---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_policy_regions"
sidebar_current: "docs-tencentcloud-datasource-teo_security_policy_regions"
description: |-
  Use this data source to query detailed information of teo securityPolicyRegions
---

# tencentcloud_teo_security_policy_regions

Use this data source to query detailed information of teo securityPolicyRegions

## Example Usage

```hcl
data "tencentcloud_teo_security_policy_regions" "securityPolicyRegions" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `geo_ip` - Region info.
  * `continent` - Name of the continent.
  * `country` - Name of the country.
  * `province` - Province of the region. Note: This field may return null, indicating that no valid value can be obtained.
  * `region_id` - Region ID.


