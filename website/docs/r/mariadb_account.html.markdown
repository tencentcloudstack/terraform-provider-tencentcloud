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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAccount](https://www.tencentcloud.com/document/product/237/16165) - This API is used to create a TencentDB account. Multiple accounts can be created for one instance. Accounts with the same username but different hosts are different accounts.
* Read: [DescribeAccounts](https://www.tencentcloud.com/document/product/237/16167) - This API is used to query the list of accounts of a specified TencentDB instance.
* Update: [ModifyAccountDescription](https://www.tencentcloud.com/document/product/237/16170) - This API is used to modify the remarks of a TencentDB account.
* Update: [ResetAccountPassword](https://www.tencentcloud.com/document/product/237/16168) - This API is used to reset the password of a TencentDB account.
* Delete: [DeleteAccount](https://www.tencentcloud.com/document/product/237/16171) - This API is used to delete a TencentDB account, which is uniquely identified by username and host.

## Import

mariadb account can be imported using the instance_id#user_name#host, e.g.
```
$ terraform import tencentcloud_mariadb_account.account tdsql-4pzs5b67#account-test#10.101.202.22
```

