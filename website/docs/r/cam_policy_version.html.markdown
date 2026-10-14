---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_policy_version"
sidebar_current: "docs-tencentcloud-resource-cam_policy_version"
description: |-
  Provides a resource to create a CAM policy version
---

# tencentcloud_cam_policy_version

Provides a resource to create a CAM policy version

## Example Usage

```hcl
resource "tencentcloud_cam_policy_version" "example" {
  policy_id      = 171173780
  set_as_default = "false"
  policy_document = jsonencode({
    "version" : "3.0",
    "statement" : [
      {
        "effect" : "allow",
        "action" : [
          "sts:AssumeRole"
        ],
        "resource" : [
          "*"
        ]
      },
      {
        "effect" : "allow",
        "action" : [
          "cos:PutObject"
        ],
        "resource" : [
          "*"
        ]
      },
      {
        "effect" : "deny",
        "action" : [
          "aa:*"
        ],
        "resource" : [
          "*"
        ]
      },
      {
        "effect" : "deny",
        "action" : [
          "aa:*"
        ],
        "resource" : [
          "*"
        ]
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

* `policy_document` - (Required, String, ForceNew) Strategic text information.
* `policy_id` - (Required, Int, ForceNew) Strategy ID.
* `set_as_default` - (Required, Bool, ForceNew) Whether to set as a version of the current strategy.
* `policy_version` - (Optional, List) Strategic version detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.

The `policy_version` object supports the following:


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.

The `policy_version` object exports the following:

* `create_date` - Strategic version creation timeNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `document` - Strategic grammar textNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `is_default_version` - Whether it is an effective version.0 means not, 1 means yesNote: This field may return NULL, indicating that the valid value cannot be obtained.
* `version_id` - Strategic version numberNote: This field may return NULL, indicating that the valid value cannot be obtained.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreatePolicyVersion](https://www.tencentcloud.com/document/product/598/36313) - This API is used to add a policy version. After creating a policy version, you can easily change the policy by changing the policy version.
* Read: [GetPolicyVersion](https://www.tencentcloud.com/document/product/598/36311) - This API is used to query policy version details.
* Delete: [DeletePolicyVersion](https://www.tencentcloud.com/document/product/598/36312) - This API is used to delete a policy version of a policy.

## Import

CAM policy version can be imported using the id, e.g.

```
terraform import tencentcloud_cam_policy_version.example 234290251#3
```

