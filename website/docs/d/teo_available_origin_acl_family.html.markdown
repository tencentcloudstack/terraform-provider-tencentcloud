---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_available_origin_acl_family"
sidebar_current: "docs-tencentcloud-datasource-teo_available_origin_acl_family"
description: |-
  Use this data source to query detailed information of TEO available origin acl family
---

# tencentcloud_teo_available_origin_acl_family

Use this data source to query detailed information of TEO available origin acl family

## Example Usage

### Query available origin acl family by zone id

```hcl
data "tencentcloud_teo_available_origin_acl_family" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
```

### Query available origin acl family by filters

```hcl
data "tencentcloud_teo_available_origin_acl_family" "example" {
  zone_id = "zone-3fkff38fyw8s"
  filters {
    name   = "OriginACLFamily"
    values = ["gaz"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Specifies the site ID.
* `filters` - (Optional, List) Filter conditions. the maximum value of Filters.Values is 20. if this parameter is left empty, all available origin acl family information under the current site will be returned. detailed filter criteria are as follows: OriginACLFamily: filter by origin acl control domain, such as gaz/mlc/emc/plat-gaz/plat-mlc/plat-emc.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered.
* `values` - (Required, Set) Value of the filtered field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `origin_acl_family_infos` - Details of origin acl family infos.
  * `active_time` - Version effective time in UTC+8, following the date and time format of the ISO 8601 standard.
  * `entire_addresses` - Origin IP range details.
    * `i_pv4` - IPv4 subnet list.
    * `i_pv6` - IPv6 subnet list.
  * `origin_acl_family` - Origin acl control domain. Values: gaz/mlc/emc/plat-gaz/plat-mlc/plat-emc.
  * `version` - Origin protection version number. Format description: standard version: gaz-xxxxx/mlc-xxxxx/emc-xxxxx; streamlined version: plat-gaz-xxxxxx/plat-mlc-xxxxxx/plat-emc-xxxxxx.


