---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_image_share_permission"
sidebar_current: "docs-tencentcloud-resource-cvm_image_share_permission"
description: |-
  Provides a resource to create a cvm image_share_permission
---

# tencentcloud_cvm_image_share_permission

Provides a resource to create a cvm image_share_permission

## Example Usage

```hcl
resource "tencentcloud_cvm_image_share_permission" "image_share_permission" {
  image_id    = "img-xxxxxx"
  account_ids = ["xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `account_ids` - (Required, Set: [`String`]) List of account IDs with which an image is shared.
* `image_id` - (Required, String, ForceNew) Image ID such as `img-gvbnzy6f`. You can only specify an image in the NORMAL state.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cvm image_share_permission can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_image_share_permission.image_share_permission image_share_permission_id
```

