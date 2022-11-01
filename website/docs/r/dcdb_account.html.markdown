---
subcategory: "dcdb"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_account"
sidebar_current: "docs-tencentcloud-resource-dcdb_account"
description: |-
  Provides a resource to create a dcdb account
---

# tencentcloud_dcdb_account

Provides a resource to create a dcdb account

## Example Usage

```hcl
resource "tencentcloud_dcdb_account" "account" {
  instance_id          = ""
  user_name            = ""
  host                 = ""
  password             = ""
  read_only            = ""
  description          = ""
  max_user_connections = ""
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String) db host.
* `instance_id` - (Required, String) instance id.
* `password` - (Required, String) password.
* `user_name` - (Required, String) account name.
* `description` - (Optional, String) description for account.
* `max_user_connections` - (Optional, Int) max user connections.
* `read_only` - (Optional, Int) whether the account is readonly. 0 means not a readonly account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dcdb account can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_account.account account_id
```

