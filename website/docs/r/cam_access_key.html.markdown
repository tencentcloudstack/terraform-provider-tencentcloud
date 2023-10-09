---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_access_key"
sidebar_current: "docs-tencentcloud-resource-cam_access_key"
description: |-
  Provides a resource to create a cam access_key
---

# tencentcloud_cam_access_key

Provides a resource to create a cam access_key

## Example Usage

```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
}
```

### Update

```hcl
resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = 100033690181
  status     = "Inactive"
}
```

## Argument Reference

The following arguments are supported:

* `access_key` - (Optional, String) Access_key is the access key identification, required when updating.
* `status` - (Optional, String) Key status, activated (Active) or inactive (Inactive), required when updating.
* `target_uin` - (Optional, Int) Specify user Uin, if not filled, the access key is created for the current user by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `secret_access_key` - Access key (key is only visible when created, please keep it properly).


## Import

cam access_key can be imported using the id, e.g.

```
terraform import tencentcloud_cam_access_key.access_key access_key_id
```

