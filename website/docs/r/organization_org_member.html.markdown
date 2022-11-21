---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_org_member"
sidebar_current: "docs-tencentcloud-resource-organization_org_member"
description: |-
  Provides a resource to create a organization org_member
---

# tencentcloud_organization_org_member

Provides a resource to create a organization org_member

## Example Usage

```hcl
resource "tencentcloud_organization_org_member" "org_member" {
  name    = "terraform_test"
  node_id = 2003721
  permission_ids = [
    1,
    2,
    3,
    4,
  ]
  policy_type = "Financial"
  remark      = "for terraform test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Member name.
* `node_id` - (Required, Int) Organization node ID.
* `permission_ids` - (Required, Set: [`Int`]) Financial management permission IDs.Valid values:- `1`: View bill.- `2`: Check balance.- `3`: Fund transfer.- `4`: Combine bill.- `5`: Issue an invoice.- `6`: Inherit discount.- `7`: Pay on behalf.value 1,2 is required.
* `policy_type` - (Required, String) Organization policy type.- `Financial`: Financial management policy.
* `pay_uin` - (Optional, String) The uin which is payment account on behalf.When `PermissionIds` contains 7, is required.
* `record_id` - (Optional, Int) Create member record ID.When create failed and needs to be recreated, is required.
* `remark` - (Optional, String) Notes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Member creation time.
* `is_allow_quit` - Whether to allow member to leave the organization.Valid values:- `Allow`.- `Denied`.
* `member_type` - Member Type.Valid values:- `Invite`: The member is invited.- `Create`: The member is created.
* `node_name` - Organization node name.
* `org_permission` - Financial management permissions.
  * `id` - Permissions ID.
  * `name` - Permissions name.
* `org_policy_name` - Organization policy name.
* `pay_name` - The member name which is payment account on behalf.
* `update_time` - Member update time.


## Import

organization org_member can be imported using the id, e.g.
```
$ terraform import tencentcloud_organization_org_member.org_member orgMember_id
```

