---
subcategory: "CDWPG"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwpg_account"
sidebar_current: "docs-tencentcloud-resource-cdwpg_account"
description: |-
  Provides a resource to create a cdwpg cdwpg_account
---

# tencentcloud_cdwpg_account

Provides a resource to create a cdwpg cdwpg_account

## Example Usage

```hcl
resource "tencentcloud_cdwpg_account" "cdwpg_account" {
  instance_id  = "cdwpg-zpiemnyd"
  user_name    = "dbadmin"
  new_password = "testpassword"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `new_password` - (Required, String) New password.
* `user_name` - (Required, String, ForceNew) Username.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cdwpg cdwpg_account can be imported using the id, e.g.

```
terraform import tencentcloud_cdwpg_account.cdwpg_account cdwpg_account_id
```

