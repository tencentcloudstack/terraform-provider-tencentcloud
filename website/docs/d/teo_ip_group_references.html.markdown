---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_ip_group_references"
sidebar_current: "docs-tencentcloud-datasource-teo_ip_group_references"
description: |-
  Use this data source to query the reference information of a specified TEO (EdgeOne) IP group.
---

# tencentcloud_teo_ip_group_references

Use this data source to query the reference information of a specified TEO (EdgeOne) IP group.

## Example Usage

```hcl
data "tencentcloud_teo_ip_group_references" "example" {
  zone_id  = "zone-3fkff38fyw8s"
  group_id = 33711
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, Int) IP group ID.
* `zone_id` - (Required, String) Site ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `references` - List of references to the IP group.
  * `entity_id` - Entity ID.
  * `entity_name` - Entity name.
  * `entity_type` - Entity type.
  * `sub_entity_id` - Sub entity ID.
  * `sub_entity_name` - Sub entity name.
  * `sub_entity_type` - Sub entity type.
  * `zone_id` - Site ID.


