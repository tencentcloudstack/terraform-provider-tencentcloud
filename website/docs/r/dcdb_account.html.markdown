---
subcategory: "TDSQL for MySQL(DCDB)"
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
  instance_id          = "tdsqlshard-kkpoxvnv"
  user_name            = "mysql"
  host                 = "127.0.0.1"
  password             = "===password==="
  read_only            = 0
  description          = "this is a test account"
  max_user_connections = 10
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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAccount](https://cloud.tencent.com/document/api/557/19996) - This API is used to create a TencentDB account. Multiple accounts can be created for one instance. Accounts with the same username but different hosts are different accounts.
* Read: [DescribeAccounts](https://cloud.tencent.com/document/api/557/19993) - This API is used to query the list of accounts of a specified TencentDB instance.
* Update: [ModifyAccountDescription](https://cloud.tencent.com/document/api/557/19984) - This API is used to modify the remarks of a TencentDB account.
* Update: [ResetAccountPassword](https://cloud.tencent.com/document/api/557/19979) - This API is used to reset the password of a TencentDB account.
* Delete: [DeleteAccount](https://cloud.tencent.com/document/api/557/19995) - This API is used to delete a TencentDB account, which is uniquely identified by username and host.
* Delete: [DescribeAccounts](https://cloud.tencent.com/document/api/557/19993) - This API is used to query the list of accounts of a specified TencentDB instance.

## Import

dcdb account can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_account.account account_id
```

