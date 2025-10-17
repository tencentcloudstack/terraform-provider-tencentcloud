---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_project_member"
sidebar_current: "docs-tencentcloud-resource-wedata_project_member"
description: |-
  Provides a resource to create a WeData project member
---

# tencentcloud_wedata_project_member

Provides a resource to create a WeData project member

~> **NOTE:** This resource must exclusive in one share unit, do not declare additional roleIds resources of this project member elsewhere.

## Example Usage

```hcl
resource "tencentcloud_wedata_project_member" "example" {
  project_id = "2983848457986924544"
  user_uin   = "100044238258"
  role_ids = [
    "308335260274237440",
    "308335260844662784"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) Project ID.
* `role_ids` - (Required, Set: [`String`]) Role ID.
* `user_uin` - (Required, String, ForceNew) User ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

WeData project member can be imported using the projectId#userUin, e.g.

```
terraform import tencentcloud_wedata_project_member.example 2983848457986924544#100044238258
```

