---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_country_area_mapping"
sidebar_current: "docs-tencentcloud-datasource-gaap_country_area_mapping"
description: |-
  Use this data source to query detailed information of gaap country area mapping
---

# tencentcloud_gaap_country_area_mapping

Use this data source to query detailed information of gaap country area mapping

## Example Usage

```hcl
data "tencentcloud_gaap_country_area_mapping" "country_area_mapping" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `country_area_mapping_list` - Country/region code mapping table.
  * `continent_inner_code` - Continental Code.
  * `continent_name` - The name of the continent.
  * `geographical_zone_inner_code` - Region code.
  * `geographical_zone_name` - Region name.
  * `nation_country_inner_code` - Country code.
  * `nation_country_name` - Country name.
  * `remark` - Annotation InformationNote: This field may return null, indicating that a valid value cannot be obtained.


