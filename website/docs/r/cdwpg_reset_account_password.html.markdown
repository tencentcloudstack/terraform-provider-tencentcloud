---
subcategory: "CDWPG"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwpg_reset_account_password"
sidebar_current: "docs-tencentcloud-resource-cdwpg_reset_account_password"
description: |-
  Provides a resource to reset cdwpg account password
---

# tencentcloud_cdwpg_reset_account_password

Provides a resource to reset cdwpg account password

## Example Usage

```hcl
resource "tencentcloud_cdwpg_reset_account_password" "cdwpg_reset_account_password" {
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

cdwpg reset account password can be imported using the id, e.g.

```
terraform import tencentcloud_cdwpg_reset_account_password.cdwpg_account cdwpg_reset_account_password_id
```

