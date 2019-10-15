---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_role_policy_attachments"
sidebar_current: "docs-tencentcloud-datasource-cam_role_policy_attachments"
description: |-
  Use this data source to query detailed information of CAM role policy attachments
---

# tencentcloud_cam_role_policy_attachments

Use this data source to query detailed information of CAM role policy attachments

## Example Usage

```hcl
data "tencentcloud_cam_role_policy_attachments" "foo" {
  role_id     = "4611686018427922725"
  policy_id   = "26800353"
  policy_type = "QCS"
  create_mode = 1
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required) Id of the attached CAM role to be queried.
* `create_mode` - (Optional) Mode of Creation of the CAM user policy attachment. 1 means the cam policy attachment is created by production, and the others indicate syntax strategy ways. 
* `policy_id` - (Optional) Id of CAM policy to be queried.
* `policy_type` - (Optional) Type of the policy strategy. 'User' means customer strategy and 'QCS' means preset strategy.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `role_policy_attachment_list` - A list of CAM role policy attachments. Each element contains the following attributes:
  * `create_mode` - Mode of Creation of the CAM role policy attachment. 1 means the cam policy attachment is created by production, and the others indicate syntax strategy ways. 
  * `create_time` - Create time of the CAM role policy attachment.
  * `policy_id` - Name of CAM role.
  * `policy_name` - Name of the policy.
  * `policy_type` - Type of the policy strategy. 'User' means customer strategy and 'QCS' means preset strategy.
  * `role_id` - Id of CAM role.


