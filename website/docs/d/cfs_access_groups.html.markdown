---
subcategory: "Cloud File Storage(CFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfs_access_groups"
sidebar_current: "docs-tencentcloud-datasource-cfs_access_groups"
description: |-
  Use this data source to query the detail information of CFS access group.
---

# tencentcloud_cfs_access_groups

Use this data source to query the detail information of CFS access group.

## Example Usage

```hcl
data "tencentcloud_cfs_access_groups" "access_groups" {
  access_group_id = "pgroup-7nx89k7l"
  name            = "test"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_id` - (Optional) A specified access group ID used to query.
* `name` - (Optional) A access group Name used to query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_group_list` - An information list of CFS access group. Each element contains the following attributes:
  * `access_group_id` - ID of the access group.
  * `create_time` - Creation time of the access group.
  * `description` - Description of the access group.
  * `name` - Name of the access group.


