---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_role_policy_attachment_by_name"
sidebar_current: "docs-tencentcloud-resource-cam_role_policy_attachment_by_name"
description: |-
  Provides a resource to create a CAM role policy attachment.
---

# tencentcloud_cam_role_policy_attachment_by_name

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

resource "tencentcloud_cam_role_policy_attachment_by_name" "role_policy_attachment_basic" {
  role_name   = var.cam_role_basic
  policy_name = var.cam_policy_basic
}
```

## Argument Reference

The following arguments are supported:

* `policy_name` - (Required, String, ForceNew) Name of the policy.
* `role_name` - (Required, String, ForceNew) Name of the attached CAM role.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_mode` - Mode of Creation of the CAM role policy attachment. `1` means the CAM policy attachment is created by production, and the others indicate syntax strategy ways.
* `create_time` - The create time of the CAM role policy attachment.
* `policy_type` - Type of the policy strategy. `User` means customer strategy and `QCS` means preset strategy.


## Import

CAM role policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_role_policy_attachment_by_name.foo ${role_name}#${policy_name}
```

