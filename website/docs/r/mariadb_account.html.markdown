---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_account"
sidebar_current: "docs-tencentcloud-resource-mariadb_account"
description: |-
  Provides a resource to create a mariadb account
---

# tencentcloud_mariadb_account

Provides a resource to create a mariadb account

## Example Usage

```hcl
resource "tencentcloud_mariadb_account" "account" {
  instance_id = "tdsql-4pzs5b67"
  user_name   = "account-test"
  host        = "10.101.202.22"
  password    = "Password123."
  read_only   = 0
  description = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String) host.
* `instance_id` - (Required, String) instance id.
* `password` - (Required, String) account password.
* `user_name` - (Required, String) user name.
* `description` - (Optional, String) account description.
* `read_only` - (Optional, Int) wether account is read only, 0 means not a read only account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mariadb account can be imported using the instance_id#user_name#host, e.g.
```
$ terraform import tencentcloud_mariadb_account.account tdsql-4pzs5b67#account-test#10.101.202.22
```

