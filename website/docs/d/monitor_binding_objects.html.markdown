---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_binding_objects"
sidebar_current: "docs-tencentcloud-datasource-monitor_binding_objects"
description: |-
  Use this data source to query policy group binding objects.
---

# tencentcloud_monitor_binding_objects

Use this data source to query policy group binding objects.

## Example Usage

```hcl
data "tencentcloud_monitor_policy_groups" "name" {
  name = "test"
}

data "tencentcloud_monitor_binding_objects" "objects" {
  group_id = data.tencentcloud_monitor_policy_groups.name.list[0].group_id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, Int) Policy group ID for query.
* `result_output_file` - (Optional, String) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list objects. Each element contains the following attributes:
  * `dimensions_json` - Represents a collection of dimensions of an object instance, json format.
  * `is_shielded` - Whether the object is shielded or not, `0` means unshielded and `1` means shielded.
  * `region` - The region where the object is located.
  * `unique_id` - Object unique ID.


