---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_role_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-cam_role_policy_attachment"
description: |-
  Provides a resource to create a CAM role policy attachment.
---

# tencentcloud_cam_role_policy_attachment

Provides a resource to create a CAM role policy attachment.

## Example Usage

```hcl
variable "cam_policy_basic" {
  default = "keep-cam-policy"
}

variable "cam_role_basic" {
  default = "keep-cam-role"
}

data "tencentcloud_cam_policies" "policy" {
  name = var.cam_policy_basic
}

data "tencentcloud_cam_roles" "roles" {
  name = var.cam_role_basic
}

resource "tencentcloud_cam_role_policy_attachment" "role_policy_attachment_basic" {
  role_id   = data.tencentcloud_cam_roles.roles.role_list.0.role_id
  policy_id = data.tencentcloud_cam_policies.policy.policy_list.0.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, ForceNew) ID of the policy.
* `role_id` - (Required, String, ForceNew) ID of the attached CAM role.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_mode` - Mode of Creation of the CAM role policy attachment. `1` means the CAM policy attachment is created by production, and the others indicate syntax strategy ways.
* `create_time` - The create time of the CAM role policy attachment.
* `policy_name` - The name of the policy.
* `policy_type` - Type of the policy strategy. `User` means customer strategy and `QCS` means preset strategy.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AttachRolePolicy](https://cloud.tencent.com/document/api/598/36226) - 绑定权限策略到角色
* Create: [ListAttachedRolePolicies](https://cloud.tencent.com/document/api/598/36220) - 获取角色绑定的策略列表
* Read: [ListAttachedRolePolicies](https://cloud.tencent.com/document/api/598/36220) - 获取角色绑定的策略列表
* Delete: [DetachRolePolicy](https://cloud.tencent.com/document/api/598/36222) - 解绑角色的策略

## Import

CAM role policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_role_policy_attachment.foo 4611686018427922725#26800353
```

