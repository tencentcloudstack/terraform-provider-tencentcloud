---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_account_groups"
sidebar_current: "docs-tencentcloud-datasource-bh_account_groups"
description: |-
  Use this data source to query detailed information of BH account groups
---

# tencentcloud_bh_account_groups

Use this data source to query detailed information of BH account groups

## Example Usage

### Query all bh account groups

```hcl
data "tencentcloud_bh_account_groups" "example" {}
```

### Query bh account groups by filter

```hcl
data "tencentcloud_bh_account_groups" "example" {
  deep_in    = 1
  parent_id  = 819729
  group_name = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `deep_in` - (Optional, Int) Whether to recursively query, 0 for non-recursive, 1 for recursive.
* `group_name` - (Optional, String) Account group name, fuzzy query.
* `page_num` - (Optional, Int) Get data from which page.
* `parent_id` - (Optional, Int) Parent account group ID, default 0, query all groups under the root account group.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `account_group_set` - Account group information.
  * `description` - Account group description.
  * `id_path` - Account group ID path.
  * `id` - Account group ID.
  * `import_type` - Account group import type.
  * `is_leaf` - Whether it is a leaf node.
  * `name_path` - Account group name path.
  * `name` - Account group name.
  * `org_id` - Source account organization ID. When using third-party import user sources, record the group ID of this group in the source organization structure.
  * `parent_id` - Parent account group ID.
  * `parent_org_id` - Parent source account organization ID. When using third-party import user sources, record the group ID of this group in the source organization structure.
  * `source` - Account group source.
  * `status` - Whether the account group has been connected, 0 means not connected, 1 means connected.
  * `user_total` - Total number of users under the account group.


