Provides a resource to create a CAM user policy attachment.

Example Usage

```hcl
variable "cam_user_basic" {
  default = "keep-cam-user"
}

resource "tencentcloud_cam_policy" "policy_basic" {
  name        = "tf_cam_attach_user_policy"
  document    =jsonencode({
    "version":"2.0",
    "statement":[
      {
        "action":["cos:*"],
        "resource":["*"],
        "effect":"allow",
      },
      {
        "effect":"allow",
        "action":["monitor:*","cam:ListUsersForGroup","cam:ListGroups","cam:GetGroup"],
        "resource":["*"],
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

Import

CAM user policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_user_policy_attachment.foo cam-test#26800353
```