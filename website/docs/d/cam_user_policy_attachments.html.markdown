---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_user_policy_attachments"
sidebar_current: "docs-tencentcloud-datasource-cam_user_policy_attachments"
description: |-
  Use this data source to query detailed information of CAM user policy attachments
---

# tencentcloud_cam_user_policy_attachments

Use this data source to query detailed information of CAM user policy attachments

## Example Usage

```hcl
# query by user_id
data "tencentcloud_cam_user_policy_attachments" "foo" {
  user_id = tencentcloud_cam_user.foo.id
}

# query by user_id and policy_id
data "tencentcloud_cam_user_policy_attachments" "bar" {
  user_id   = tencentcloud_cam_user.foo.id
  policy_id = tencentcloud_cam_policy.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required) ID of the attached CAM user to be queried.
* `create_mode` - (Optional) Mode of Creation of the CAM user policy attachment. `1` means the CAM policy attachment is created by production, and the others indicate syntax strategy ways.
* `policy_id` - (Optional) ID of CAM policy to be queried.
* `policy_type` - (Optional) Type of the policy strategy. 'User' means customer strategy and 'QCS' means preset strategy.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_policy_attachment_list` - A list of CAM user policy attachments. Each element contains the following attributes:
  * `create_mode` - Mode of Creation of the CAM user policy attachment. `1` means the cam policy attachment is created by production, and the others indicate syntax strategy ways.
  * `create_time` - The create time of the CAM user policy attachment.
  * `policy_id` - Name of CAM user.
  * `policy_name` - The name of the policy.
  * `policy_type` - Type of the policy strategy. 'User' means customer strategy and 'QCS' means preset strategy.
  * `user_id` - ID of CAM user.


