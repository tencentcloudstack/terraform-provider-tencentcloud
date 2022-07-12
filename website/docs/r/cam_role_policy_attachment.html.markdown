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
resource "tencentcloud_cam_role_policy_attachment" "foo" {
  role_id   = tencentcloud_cam_role.foo.id
  policy_id = tencentcloud_cam_policy.foo.id
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


## Import

CAM role policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_role_policy_attachment.foo 4611686018427922725#26800353
```

