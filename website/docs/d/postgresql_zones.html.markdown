---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_zones"
sidebar_current: "docs-tencentcloud-datasource-postgresql_zones"
description: |-
  Use this data source to query detailed information of postgresql zones
---

# tencentcloud_postgresql_zones

Use this data source to query detailed information of postgresql zones

## Example Usage

```hcl
data "tencentcloud_postgresql_zones" "zones" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `zone_set` - AZ information set.
  * `standby_zone_set` - AZs that can be used as standby when this AZ is primaryNote: this field may return `null`, indicating that no valid values can be obtained.
  * `zone_id` - AZ number.
  * `zone_name` - AZ name.
  * `zone_state` - Availability status. Valid values:`UNAVAILABLE`.`AVAILABLE`.`SELLOUT`.`SUPPORTMODIFYONLY` (supports configuration adjustment).
  * `zone_support_ipv6` - Whether the AZ supports IPv6 address access.
  * `zone` - AZ abbreviation.


