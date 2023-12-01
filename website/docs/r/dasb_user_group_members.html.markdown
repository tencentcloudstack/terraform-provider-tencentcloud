---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_user_group_members"
sidebar_current: "docs-tencentcloud-resource-dasb_user_group_members"
description: |-
  Provides a resource to create a dasb user_group_members
---

# tencentcloud_dasb_user_group_members

Provides a resource to create a dasb user_group_members

## Example Usage

```hcl
resource "tencentcloud_dasb_user_group_members" "example" {
  user_group_id = 3
  member_id_set = [1, 2, 3]
}
```

## Argument Reference

The following arguments are supported:

* `member_id_set` - (Required, Set: [`Int`], ForceNew) Collection of member user IDs.
* `user_group_id` - (Required, Int, ForceNew) User Group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dasb user_group_members can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_user_group_members.example 3#1,2,3
```

