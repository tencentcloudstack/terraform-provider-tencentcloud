---
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
  user_id   = "cam-test"
  policy_id = "26800353"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew) Id of the policy.
* `user_id` - (Required, ForceNew) Id of the attached CAM user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_mode` - Mode of Creation of the CAM user policy attachment. 1 means the CAM policy attachment is created by production, and the others indicate syntax strategy ways.
* `create_time` - Create time of the CAM user policy attachment.
* `policy_name` - Name of the policy.
* `policy_type` - Type of the policy strategy. 'User' means customer strategy and 'QCS' means preset strategy.


## Import

CAM user policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_user_policy_attachment.foo cam-test#26800353
```

