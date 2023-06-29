---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_account"
sidebar_current: "docs-tencentcloud-resource-mysql_account"
description: |-
  Provides a MySQL account resource for database management. A MySQL instance supports multiple database account.
---

# tencentcloud_mysql_account

Provides a MySQL account resource for database management. A MySQL instance supports multiple database account.

## Example Usage

```hcl
resource "tencentcloud_mysql_account" "default" {
  mysql_id             = "terraform-test-local-database"
  name                 = "tf_test"
  password             = "********"
  description          = "My test account"
  max_user_connections = 10
}
```

## Argument Reference

The following arguments are supported:

* `mysql_id` - (Required, String, ForceNew) Instance ID to which the account belongs.
* `name` - (Required, String, ForceNew) Account name.
* `password` - (Required, String) Operation password.
* `description` - (Optional, String) Database description.
* `host` - (Optional, String) Account host, default is `%`.
* `max_user_connections` - (Optional, Int) The maximum number of available connections for a new account, the default value is 10240, and the maximum value that can be set is 10240.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql account can be imported using the mysqlId#accountName, e.g.

```
terraform import tencentcloud_mysql_account.default cdb-gqg6j82x#tf_account
```

