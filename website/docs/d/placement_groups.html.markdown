---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_placement_groups"
sidebar_current: "docs-tencentcloud-datasource-placement_groups"
description: |-
  Use this data source to query placement groups.
---

# tencentcloud_placement_groups

Use this data source to query placement groups.

## Example Usage

```hcl
data "tencentcloud_placement_groups" "foo" {
  placement_group_id = "ps-21q9ibvr"
  name               = "test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the placement group to be queried.
* `placement_group_id` - (Optional) ID of the placement group to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `placement_group_list` - An information list of placement group. Each element contains the following attributes:
  * `create_time` - Creation time of the placement group.
  * `current_num` - Number of hosts in the placement group.
  * `cvm_quota_total` - Maximum number of hosts in the placement group.
  * `instance_ids` - Host IDs in the placement group.
  * `name` - Name of the placement group.
  * `placement_group_id` - ID of the placement group.
  * `type` - Type of the placement group.


