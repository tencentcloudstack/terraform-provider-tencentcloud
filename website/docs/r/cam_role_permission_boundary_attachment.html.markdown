---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_role_permission_boundary_attachment"
sidebar_current: "docs-tencentcloud-resource-cam_role_permission_boundary_attachment"
description: |-
  Provides a resource to create a cam role_permission_boundary_attachment
---

# tencentcloud_cam_role_permission_boundary_attachment

Provides a resource to create a cam role_permission_boundary_attachment

## Example Usage

```hcl
resource "tencentcloud_cam_role_permission_boundary_attachment" "role_permission_boundary_attachment" {
  policy_id = 1
  role_name = "test-cam-tag"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, Int, ForceNew) Role ID.
* `role_id` - (Optional, String, ForceNew) Role ID (at least one should be filled in with the role name).
* `role_name` - (Optional, String, ForceNew) Role name (at least one should be filled in with the role ID).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cam role_permission_boundary_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cam_role_permission_boundary_attachment.role_permission_boundary_attachment role_permission_boundary_attachment_id
```

