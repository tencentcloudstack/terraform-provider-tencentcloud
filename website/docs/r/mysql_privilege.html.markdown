---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_privilege"
sidebar_current: "docs-tencentcloud-resource-mysql_privilege"
description: |-
  Provides a mysql account privilege resource to grant different access privilege to different database. A database can be granted by multiple account.
---

# tencentcloud_mysql_privilege

Provides a mysql account privilege resource to grant different access privilege to different database. A database can be granted by multiple account.

## Example Usage

```hcl
resource "tencentcloud_mysql_instance" "default" {
  mem_size          = 1000
  volume_size       = 25
  instance_name     = "guagua"
  engine_version    = "5.7"
  root_password     = "0153Y474"
  availability_zone = "ap-guangzhou-3"
  internet_service  = 1

}

resource "tencentcloud_mysql_account" "mysql_account2" {
  mysql_id    = tencentcloud_mysql_instance.default.id
  name        = "test11"
  password    = "test1234"
  description = "test from terraform"
}

resource "tencentcloud_mysql_privilege" "tttt" {
  mysql_id     = tencentcloud_mysql_instance.default.id
  account_name = tencentcloud_mysql_account.mysql_account2.name
  global       = ["TRIGGER"]
  database {
    privileges    = ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE"]
    database_name = "sys"
  }
  database {
    privileges    = ["SELECT"]
    database_name = "performance_schema"
  }

  table {
    privileges    = ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE"]
    database_name = "mysql"
    table_name    = "slow_log"
  }

  table {
    privileges    = ["SELECT", "INSERT", "UPDATE"]
    database_name = "mysql"
    table_name    = "user"
  }

  column {
    privileges    = ["SELECT", "INSERT", "UPDATE", "REFERENCES"]
    database_name = "mysql"
    table_name    = "user"
    column_name   = "host"
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, ForceNew) Account name.the forbidden value is:root,mysql.sys,tencentroot.
* `global` - (Required) Global privileges. available values for Privileges:SELECT,INSERT,UPDATE,DELETE,CREATE,PROCESS,DROP,REFERENCES,INDEX,ALTER,SHOW DATABASES,CREATE TEMPORARY TABLES,LOCK TABLES,EXECUTE,CREATE VIEW,SHOW VIEW,CREATE ROUTINE,ALTER ROUTINE,EVENT,TRIGGER.
* `mysql_id` - (Required, ForceNew) Instance ID.
* `column` - (Optional) Column privileges list.
* `database` - (Optional) Database privileges list.
* `table` - (Optional) Table privileges list.

The `column` object supports the following:

* `column_name` - (Required) Column name.
* `database_name` - (Required) Database name.
* `privileges` - (Required) Column privilege.available values for Privileges:SELECT,INSERT,UPDATE,REFERENCES.
* `table_name` - (Required) Table name.

The `database` object supports the following:

* `database_name` - (Required) Database name.
* `privileges` - (Required) Database privilege.available values for Privileges:SELECT,INSERT,UPDATE,DELETE,CREATE,DROP,REFERENCES,INDEX,ALTER,CREATE TEMPORARY TABLES,LOCK TABLES,EXECUTE,CREATE VIEW,SHOW VIEW,CREATE ROUTINE,ALTER ROUTINE,EVENT,TRIGGER.

The `table` object supports the following:

* `database_name` - (Required) Database name.
* `privileges` - (Required) Table privilege.available values for Privileges:SELECT,INSERT,UPDATE,DELETE,CREATE,DROP,REFERENCES,INDEX,ALTER,CREATE VIEW,SHOW VIEW,TRIGGER.
* `table_name` - (Required) Table name.


