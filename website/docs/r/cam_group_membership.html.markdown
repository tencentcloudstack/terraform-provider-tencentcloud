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
variable "cam_group_basic" {
  default = "keep-cam-group"
}

data "tencentcloud_cam_groups" "groups" {
  name = var.cam_group_basic
}

resource "tencentcloud_cam_user" "foo" {
  name                = "tf_cam_user"
  remark              = "tf_user_remark"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "12345678910"
  country_code        = "86"
  email               = "1234@qq.com"
  force_delete        = true
}

resource "tencentcloud_cam_group_membership" "group_membership_basic" {
  group_id   = data.tencentcloud_cam_groups.groups.group_list.0.group_id
  user_names = [tencentcloud_cam_user.foo.id]
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



## Import

CAM group membership can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group_membership.foo 12515263
```

