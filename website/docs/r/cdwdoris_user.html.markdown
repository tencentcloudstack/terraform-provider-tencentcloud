---
subcategory: "CdwDoris"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwdoris_user"
sidebar_current: "docs-tencentcloud-resource-cdwdoris_user"
description: |-
  Provides a resource to create a cdwdoris cdwdoris_user
---

# tencentcloud_cdwdoris_user

Provides a resource to create a cdwdoris cdwdoris_user

## Example Usage

```hcl
resource "tencentcloud_cdwdoris_user" "cdwdoris_user" {
  user_info = {
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_type` - (Required, String, ForceNew) Api type.
* `user_info` - (Required, List) User info.
* `user_privilege` - (Optional, Int, ForceNew) User permission type. 0: Ordinary user; 1: Administrator.

The `user_info` object supports the following:

* `instance_id` - (Required, String) Instance ID.
* `password` - (Required, String) Password.
* `username` - (Required, String) User name.
* `cam_ranger_group_ids` - (Optional, List) Ranger group id list.
* `cam_uin` - (Optional, String) The bound sub user uin.
* `describe` - (Optional, String) Describe.
* `white_host` - (Optional, String) The IP the user linked from.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cdwdoris cdwdoris_user can be imported using the id, e.g.

```
terraform import tencentcloud_cdwdoris_user.cdwdoris_user cdwdoris_user_id
```

