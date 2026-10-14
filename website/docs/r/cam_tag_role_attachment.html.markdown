---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_tag_role_attachment"
sidebar_current: "docs-tencentcloud-resource-cam_tag_role_attachment"
description: |-
  Provides a resource to create a CAM tag role
---

# tencentcloud_cam_tag_role_attachment

Provides a resource to create a CAM tag role

## Example Usage

### Create by role_id

```hcl
resource "tencentcloud_cam_tag_role_attachment" "example" {
  role_id = "4611686018441060141"

  tags {
    key   = "tagKey"
    value = "tagValue"
  }
}
```

### Create by role_name

```hcl
resource "tencentcloud_cam_tag_role_attachment" "example" {
  role_name = "tf-example"

  tags {
    key   = "tagKey"
    value = "tagValue"
  }
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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [TagRole](https://cloud.tencent.com/document/api/598/74379) - This API is used to bind tags to a role.
* Read: [GetRole](https://cloud.tencent.com/document/api/598/36221) - This API (GetRole) is used to get the details of a specified role.
* Delete: [UntagRole](https://cloud.tencent.com/document/api/598/74378) - This API is used to unbind tags from a role.

## Import

CAM tag role can be imported using the id, e.g.

```
# Please use role_name#role_id
terraform import tencentcloud_cam_tag_role_attachment.example tf-example#4611686018441060141
```

