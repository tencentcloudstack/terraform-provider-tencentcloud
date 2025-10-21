---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_groups"
sidebar_current: "docs-tencentcloud-datasource-identity_center_groups"
description: |-
  Use this data source to query detailed information of identity center groups
---

# tencentcloud_identity_center_groups

Use this data source to query detailed information of identity center groups

## Example Usage

```hcl
data "tencentcloud_identity_center_groups" "identity_center_groups" {
  zone_id = "z-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Space ID.
* `filter_users` - (Optional, Set: [`String`]) Filtered user. IsSelected=1 will be returned for the user group associated with this user.
* `filter` - (Optional, String) Filter criterion. Format: <Attribute> <Operator> <Value>, case-insensitive. Currently, <Attribute> supports only GroupName, and <Operator> supports only eq (Equals) and sw (Start With). For example, Filter = "GroupName sw test" indicates querying all user groups with names starting with test; Filter = "GroupName eq testgroup" indicates querying the user group with the name testgroup.
* `group_type` - (Optional, String) User group type. Manual: manually created; Synchronized: externally imported.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_field` - (Optional, String) Sorting field, which currently only supports CreateTime. The default is the CreateTime field.
* `sort_type` - (Optional, String) Sorting type. Desc: descending order; Asc: ascending order. It should be set along with SortField.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `groups` - User group list.


