---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_modify_user_typ_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_modify_user_typ_operation"
description: |-
  Provides a resource to create a DLC modify user typ operation
---

# tencentcloud_dlc_modify_user_typ_operation

Provides a resource to create a DLC modify user typ operation

## Example Usage

```hcl
resource "tencentcloud_dlc_modify_user_typ_operation" "example" {
  user_id   = "127382378"
  user_type = "ADMIN"
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String, ForceNew) User ID.
* `user_type` - (Required, String, ForceNew) Types that users modify. ADMIN: administrators; COMMON: general users.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



