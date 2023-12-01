---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_tag_role_attachment"
sidebar_current: "docs-tencentcloud-resource-cam_tag_role_attachment"
description: |-
  Provides a resource to create a cam tag_role
---

# tencentcloud_cam_tag_role_attachment

Provides a resource to create a cam tag_role

## Example Usage

```hcl
resource "tencentcloud_cam_tag_role_attachment" "tag_role" {
  tags {
    key   = "test1"
    value = "test1"
  }
  role_id = "test-cam-tag"
}
```

## Argument Reference

The following arguments are supported:

* `tags` - (Required, List, ForceNew) Label.
* `role_id` - (Optional, String, ForceNew) Character ID, at least one input with the character name.
* `role_name` - (Optional, String, ForceNew) Character name, at least one input with the character ID.

The `tags` object supports the following:

* `key` - (Required, String) Label.
* `value` - (Required, String) Label.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cam tag_role can be imported using the id, e.g.

```
terraform import tencentcloud_cam_tag_role_attachment.tag_role tag_role_id
```

