---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_user_group_attachment"
sidebar_current: "docs-tencentcloud-resource-identity_center_user_group_attachment"
description: |-
  Provides a resource to create an identity center user group attachment
---

# tencentcloud_identity_center_user_group_attachment

Provides a resource to create an identity center user group attachment

## Example Usage

```hcl
resource "tencentcloud_identity_center_user_group_attachment" "identity_center_user_group_attachment" {
  zone_id  = "z-xxxxxx"
  user_id  = "u-xxxxxx"
  group_id = "g-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, ForceNew) User group ID.
* `user_id` - (Required, String, ForceNew) User ID.
* `zone_id` - (Required, String, ForceNew) Zone id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddUserToGroup](https://cloud.tencent.com/document/api/850/109921) - 为用户组添加用户
* Read: [ListJoinedGroupsForUser](https://cloud.tencent.com/document/api/850/109915) - 查询用户加入的用户组
* Delete: [RemoveUserFromGroup](https://cloud.tencent.com/document/api/850/109914) - 从用户组中移除用户

## Import

organization identity_center_user_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_user_group_attachment.identity_center_user_group_attachment ${zoneId}#${groupId}#${userId}
```

