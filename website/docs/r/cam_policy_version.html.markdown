---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_policy_version"
sidebar_current: "docs-tencentcloud-resource-cam_policy_version"
description: |-
  Provides a resource to create a cam policy_version
---

# tencentcloud_cam_policy_version

Provides a resource to create a cam policy_version

## Example Usage

```hcl
resource "tencentcloud_cam_policy_version" "policy_version" {
  policy_id = 171173780
  policy_document = jsonencode({
    "version" : "2.0",
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
  set_as_default = "false"
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



## Import

cam policy_version can be imported using the id, e.g.

```
terraform import tencentcloud_cam_policy_version.policy_version policy_version_id
```

