---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_available_origin_acl_family"
sidebar_current: "docs-tencentcloud-datasource-teo_available_origin_acl_family"
description: |-
  Use this data source to query available origin ACL families of TEO origin protection
---

# tencentcloud_teo_available_origin_acl_family

Use this data source to query available origin ACL families of TEO origin protection

## Example Usage

### Query all available origin ACL families by zone id

```hcl
data "tencentcloud_teo_available_origin_acl_family" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
```

### Query available origin ACL families by filter

```hcl
data "tencentcloud_teo_available_origin_acl_family" "example" {
  zone_id = "zone-3fkff38fyw8s"

  filters {
    name   = "OriginACLFamily"
    values = ["gaz", "mlc"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Site ID.
* `filters` - (Optional, List) Filter conditions. The upper limit of Filters.Values is 20. Only `OriginACLFamily` is supported, used to filter by origin ACL control domain.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter name. Valid value: `OriginACLFamily`.
* `values` - (Required, Set) Filter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `origin_acl_family_info_set` - Origin ACL family info list.
  * `active_time` - Version active time, Beijing time UTC+8, following ISO 8601 standard.
  * `entire_addresses` - Origin IP CIDR details.
    * `ipv4` - IPv4 CIDR list.
    * `ipv6` - IPv6 CIDR list.
  * `origin_acl_family` - Origin ACL control domain.
  * `version` - Origin protection version number.
* `total_count` - Total number of available origin ACL families.


