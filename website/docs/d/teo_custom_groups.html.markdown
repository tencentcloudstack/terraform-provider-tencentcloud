---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_custom_groups"
sidebar_current: "docs-tencentcloud-data_source-teo_custom_groups"
description: |-
  Use this data source to query the list of EdgeOne (TEO) custom groups.
---

# tencentcloud_teo_custom_groups

Use this data source to query the list of EdgeOne (TEO) custom groups.

## Example Usage

### Query all custom groups of a site

```hcl
data "tencentcloud_teo_custom_groups" "example" {
  zone_id = "zone-27h0vbm5w1e"
}
```

### Query custom groups filtered by group type

```hcl
data "tencentcloud_teo_custom_groups" "example" {
  zone_id = "zone-27h0vbm5w1e"

  group_types = ["TLSFingerprint", "UA"]
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Site ID.
* `group_types` - (Optional, List) List of custom group types. If not specified, all types of custom groups are queried. Valid values: `TLSFingerprint`, `UA`, `Referer`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

    * `content` - Entry content.
    * `item_type` - Entry type. Valid values: `JA3` (only available when the group type is TLSFingerprint), `JA4` (only available when the group type is TLSFingerprint), `UA` (only available when the group type is UA), `Referer` (only available when the group type is Referer).
  * `group_id` - Unique ID of the custom group.
  * `group_item_total_count` - Total number of entries in the custom group.
  * `group_items` - List of custom group entries.
  * `group_name` - Name of the custom group.
  * `group_type` - Type of the custom group. Valid values: `TLSFingerprint`, `UA`, `Referer`.
  * `ref_count` - Number of times the custom group is referenced.
* `custom_groups` - List of custom groups that match the group type conditions.
* `group_total_count` - Total number of custom groups that match the group type conditions.


