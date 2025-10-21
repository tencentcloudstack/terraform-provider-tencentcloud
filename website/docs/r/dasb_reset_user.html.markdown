---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_reset_user"
sidebar_current: "docs-tencentcloud-resource-dasb_reset_user"
description: |-
  Provides a resource to create a dasb reset_user
---

# tencentcloud_dasb_reset_user

Provides a resource to create a dasb reset_user

## Example Usage

```hcl
resource "tencentcloud_dasb_reset_user" "example" {
  user_id = 16
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, Int, ForceNew) User Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



