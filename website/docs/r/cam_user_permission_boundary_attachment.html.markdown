---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_user_permission_boundary_attachment"
sidebar_current: "docs-tencentcloud-resource-cam_user_permission_boundary_attachment"
description: |-
  Provides a resource to create a CAM user permission boundary
---

# tencentcloud_cam_user_permission_boundary_attachment

Provides a resource to create a CAM user permission boundary

## Example Usage

```hcl
resource "tencentcloud_cam_user_permission_boundary_attachment" "example" {
  target_uin = 100037718101
  policy_id  = 234290251
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, Int, ForceNew) Policy ID.
* `target_uin` - (Required, Int, ForceNew) Sub account Uin.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [PutUserPermissionsBoundary](https://www.tencentcloud.com/document/api/598/37925)
* Read: [GetUserPermissionBoundary](https://www.tencentcloud.com/search/?wd=GetUserPermissionBoundary&type=document&lang=en)
* Delete: [DeleteUserPermissionsBoundary](https://www.tencentcloud.com/document/api/598/37926)

## Import

CAM user permission boundary can be imported using the id, e.g.

```
terraform import tencentcloud_cam_user_permission_boundary_attachment.example 100037718101#234290251
```

