---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_user"
sidebar_current: "docs-tencentcloud-resource-dlc_user"
description: |-
  Provides a resource to create a dlc user
---

# tencentcloud_dlc_user

Provides a resource to create a dlc user

## Example Usage

```hcl
resource "tencentcloud_dlc_user" "user" {
  user_id          = "100027012454"
  user_type        = "COMMON"
  user_alias       = "terraform-test"
  user_description = "for terraform test"
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String, ForceNew) The sub-user uin that needs to be authorized.
* `user_alias` - (Optional, String) User alias, the character length is less than 50.
* `user_description` - (Optional, String) User description information, easy to distinguish between different users.
* `user_type` - (Optional, String) User Type. `ADMIN` or `COMMONN`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `work_group_ids` - A collection of workgroup IDs bound to the user.


## Import

dlc user can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_user.user user_id
```

