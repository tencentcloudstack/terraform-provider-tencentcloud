---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_account"
sidebar_current: "docs-tencentcloud-tencentcloud_mysql_account"
description: |-
Provides a MySQL account resource for database management. A MySQL instance supports multiple database account.

---
#tencentcloud_mysql_account

Provides a MySQL account resource for database management. A MySQL instance supports multiple database account.


##Example Usage

```
resource "tencentcloud_mysql_account" "default" {
  mysql_id = "my-test-database"
  name = "tf_account"
  password = "..."
  description = "My test account"
}

```

##Argument Reference


The following arguments are supported:

- `mysql_id` - (Required) Instance ID to which the account belongs.

- `name` - (Required) Account name.

- `password` - (Required) Operation password.

- `description` - (Optional) Database description.