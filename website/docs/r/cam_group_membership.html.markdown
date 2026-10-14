---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_group_membership"
sidebar_current: "docs-tencentcloud-resource-cam_group_membership"
description: |-
  Provides a resource to create a CAM group membership.
---

# tencentcloud_cam_group_membership

Provides a resource to create a CAM group membership.

## Example Usage

```hcl
// create cam user
resource "tencentcloud_cam_user" "example" {
  name                = "tf-example"
  remark              = "remark."
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Password@123"
  phone_num           = "18611111111"
  email               = "example@tencent.com"
  country_code        = "86"
  force_delete        = true
  tags = {
    createBy = "Terraform"
  }
}

// create cam group
resource "tencentcloud_cam_group" "example" {
  name   = "tf-example"
  remark = "remark."
}

// create cam group membership
resource "tencentcloud_cam_group_membership" "example" {
  group_id   = tencentcloud_cam_group.example.id
  user_names = [tencentcloud_cam_user.example.id]
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) ID of CAM group.
* `user_ids` - (Optional, Set: [`String`], **Deprecated**) It has been deprecated from version 1.59.5. Use `user_names` instead. ID set of the CAM group members.
* `user_names` - (Optional, Set: [`String`]) User name set as ID of the CAM group members.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddUserToGroup](https://cloud.tencent.com/document/api/598/34594) - This API is used to add users to a user group.
* Create: [ListUsersForGroup](https://cloud.tencent.com/document/api/598/34586) - This API is used to query the list of users associated with a user group.
* Read: [ListUsersForGroup](https://cloud.tencent.com/document/api/598/34586) - This API is used to query the list of users associated with a user group.
* Delete: [RemoveUserFromGroup](https://cloud.tencent.com/document/api/598/34585) - This API is used to delete users from a user group.

## Import

CAM group membership can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group_membership.example 353251
```

