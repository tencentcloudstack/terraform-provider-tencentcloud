---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_group_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-cam_group_policy_attachment"
description: |-
  Provides a resource to create a CAM group policy attachment.
---

# tencentcloud_cam_group_policy_attachment

Provides a resource to create a CAM group policy attachment.

## Example Usage

```hcl
variable "cam_policy_basic" {
  default = "keep-cam-policy"
}

variable "cam_group_basic" {
  default = "keep-cam-group"
}

data "tencentcloud_cam_groups" "groups" {
  name = var.cam_group_basic
}

data "tencentcloud_cam_policies" "policy" {
  name = var.cam_policy_basic
}

resource "tencentcloud_cam_group_policy_attachment" "group_policy_attachment_basic" {
  group_id  = data.tencentcloud_cam_groups.groups.group_list.0.group_id
  policy_id = data.tencentcloud_cam_policies.policy.policy_list.0.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, ForceNew) ID of the attached CAM group.
* `policy_id` - (Required, String, ForceNew) ID of the policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_mode` - Mode of Creation of the CAM group policy attachment. `1` means the cam policy attachment is created by production, and the others indicate syntax strategy ways.
* `create_time` - Create time of the CAM group policy attachment.
* `policy_name` - Name of the policy.
* `policy_type` - Type of the policy strategy. 'Group' means customer strategy and 'QCS' means preset strategy.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AttachGroupPolicy](https://www.tencentcloud.com/document/product/598/32250) - This API (AttachGroupPolicy) is used to associate a policy with a user group.
* Create: [ListAttachedGroupPolicies](https://www.tencentcloud.com/document/product/598/32243) - This API (ListAttachedGroupPolicies) is used to query the list of policies associated with a user group.
* Read: [ListAttachedGroupPolicies](https://www.tencentcloud.com/document/product/598/32243) - This API (ListAttachedGroupPolicies) is used to query the list of policies associated with a user group.
* Delete: [DetachGroupPolicy](https://www.tencentcloud.com/document/product/598/32246) - This API (DetachGroupPolicy) is used to unassociate a policy and a user group.

## Import

CAM group policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group_policy_attachment.foo 12515263#26800353
```

