---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_user_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-cam_user_policy_attachment"
description: |-
  Provides a resource to create a CAM user policy attachment.
---

# tencentcloud_cam_user_policy_attachment

Provides a resource to create a CAM user policy attachment.

## Example Usage

```hcl
variable "cam_user_basic" {
  default = "keep-cam-user"
}

resource "tencentcloud_cam_policy" "policy_basic" {
  name = "tf_cam_attach_user_policy"
  document = jsonencode({
    "version" : "2.0",
    "statement" : [
      {
        "action" : ["cos:*"],
        "resource" : ["*"],
        "effect" : "allow",
      },
      {
        "effect" : "allow",
        "action" : ["monitor:*", "cam:ListUsersForGroup", "cam:ListGroups", "cam:GetGroup"],
        "resource" : ["*"],
      }
    ]
  })
  description = "tf_test"
}

data "tencentcloud_cam_users" "users" {
  name = var.cam_user_basic
}

resource "tencentcloud_cam_user_policy_attachment" "user_policy_attachment_basic" {
  user_name = data.tencentcloud_cam_users.users.user_list.0.user_id
  policy_id = tencentcloud_cam_policy.policy_basic.id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, ForceNew) ID of the policy.
* `user_id` - (Optional, String, ForceNew, **Deprecated**) It has been deprecated from version 1.59.5. Use `user_name` instead. ID of the attached CAM user.
* `user_name` - (Optional, String, ForceNew) Name of the attached CAM user as uniq key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_mode` - Mode of Creation of the CAM user policy attachment. `1` means the CAM policy attachment is created by production, and the others indicate syntax strategy ways.
* `create_time` - Create time of the CAM user policy attachment.
* `policy_name` - Name of the policy.
* `policy_type` - Type of the policy strategy. `User` means customer strategy and `QCS` means preset strategy.


## Import

CAM user policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_user_policy_attachment.foo cam-test#26800353
```

