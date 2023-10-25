---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_modify_user_typ_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_modify_user_typ_operation"
description: |-
  Provides a resource to create a dlc modify_user_typ_operation
---

# tencentcloud_dlc_modify_user_typ_operation

Provides a resource to create a dlc modify_user_typ_operation

## Example Usage

```hcl
resource "tencentcloud_dlc_modify_user_typ_operation" "modify_user_typ_operation" {
  user_id   = "127382378"
  user_type = "ADMIN"
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String, ForceNew) User id (uin), if left blank, it defaults to the caller's sub-uin.
* `user_type` - (Required, String, ForceNew) User type, only support: ADMIN: ddministrator/COMMON: ordinary user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dlc modify_user_typ_operation can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_modify_user_typ_operation.modify_user_typ_operation modify_user_typ_operation_id
```

