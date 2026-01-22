---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_user"
sidebar_current: "docs-tencentcloud-resource-dlc_user"
description: |-
  Provides a resource to create a DLC user
---

# tencentcloud_dlc_user

Provides a resource to create a DLC user

## Example Usage

```hcl
resource "tencentcloud_dlc_user" "example" {
  user_id          = "100027012454"
  user_type        = "COMMON"
  user_alias       = "terraform-test"
  user_description = "for terraform test"
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String, ForceNew) Sub-user UIN that needs to be granted permissions. It can be checked through the upper right corner of Tencent Cloud Console -> Account Information -> Account ID.
* `user_alias` - (Optional, String) User alias, and its characters are less than 50.
* `user_description` - (Optional, String) User description, which can make it easy to identify different users.
* `user_type` - (Optional, String) Types of users. ADMIN: administrators; COMMON: general users. When the type of user is administrator, the collections of permissions and bound working groups cannot be set. Administrators own all the permissions by default. If the parameter is not filled in, it will be COMMON by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `work_group_ids` - Collection of IDs of working groups bound to users.


## Import

dlc user can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_user.example 100027012454
```

