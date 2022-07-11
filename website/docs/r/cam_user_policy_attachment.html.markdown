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
resource "tencentcloud_cam_user_policy_attachment" "foo" {
  user_id   = tencentcloud_cam_user.foo.id
  policy_id = tencentcloud_cam_policy.foo.id
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

