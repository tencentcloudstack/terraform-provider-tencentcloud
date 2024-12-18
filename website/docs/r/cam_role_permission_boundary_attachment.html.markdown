---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_role_permission_boundary_attachment"
sidebar_current: "docs-tencentcloud-resource-cam_role_permission_boundary_attachment"
description: |-
  Provides a resource to create a CAM role permission boundary attachment
---

# tencentcloud_cam_role_permission_boundary_attachment

Provides a resource to create a CAM role permission boundary attachment

## Example Usage

### Use role_name

```hcl
resource "tencentcloud_cam_role_permission_boundary_attachment" "example" {
  policy_id = 1
  role_name = "tf-example"
}
```

### Use role_id

```hcl
resource "tencentcloud_cam_role_permission_boundary_attachment" "example" {
  policy_id = 1
  role_id   = "4611686018441060141"
}
```

### Use all

```hcl
resource "tencentcloud_cam_role_permission_boundary_attachment" "example" {
  policy_id = 1
  role_name = "tf-example"
  role_id   = "4611686018441060141"
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

CAM role permission boundary attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cam_role_permission_boundary_attachment.example 1##tf-example

terraform import tencentcloud_cam_role_permission_boundary_attachment.example 1#4611686018441060141#

terraform import tencentcloud_cam_role_permission_boundary_attachment.example 1#4611686018441060141#tf-example
```

