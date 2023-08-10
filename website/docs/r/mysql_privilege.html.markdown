---
subcategory: "TencentDB for MySQL(cdb)"
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
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-mysql"
  description = "mysql test"
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  slave_sync_mode   = 1
  instance_name     = "tf-example-mysql"
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [tencentcloud_security_group.security_group.id]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_mysql_account" "example" {
  mysql_id             = tencentcloud_mysql_instance.example.id
  name                 = "tf_example"
  password             = "Qwer@234"
  description          = "desc."
  max_user_connections = 10
}

resource "tencentcloud_mysql_privilege" "example" {
  mysql_id     = tencentcloud_mysql_instance.example.id
  account_name = tencentcloud_mysql_account.example.name
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

* `account_name` - (Required, String, ForceNew) Account name.the forbidden value is:root,mysql.sys,tencentroot.
* `global` - (Required, Set: [`String`]) Global privileges. available values for Privileges:ALTER,ALTER ROUTINE,CREATE,CREATE ROUTINE,CREATE TEMPORARY TABLES,CREATE USER,CREATE VIEW,DELETE,DROP,EVENT,EXECUTE,INDEX,INSERT,LOCK TABLES,PROCESS,REFERENCES,RELOAD,REPLICATION CLIENT,REPLICATION SLAVE,SELECT,SHOW DATABASES,SHOW VIEW,TRIGGER,UPDATE.
* `mysql_id` - (Required, String, ForceNew) Instance ID.
* `account_host` - (Optional, String, ForceNew) Account host, default is `%`.
* `column` - (Optional, Set) Column privileges list.
* `database` - (Optional, Set) Database privileges list.
* `table` - (Optional, Set) Table privileges list.

The `column` object supports the following:

* `column_name` - (Required, String) Column name.
* `database_name` - (Required, String) Database name.
* `privileges` - (Required, Set) Column privilege.available values for Privileges:SELECT,INSERT,UPDATE,REFERENCES.
* `table_name` - (Required, String) Table name.

The `database` object supports the following:

* `database_name` - (Required, String) Database name.
* `privileges` - (Required, Set) Database privilege.available values for Privileges:SELECT,INSERT,UPDATE,DELETE,CREATE,DROP,REFERENCES,INDEX,ALTER,CREATE TEMPORARY TABLES,LOCK TABLES,EXECUTE,CREATE VIEW,SHOW VIEW,CREATE ROUTINE,ALTER ROUTINE,EVENT,TRIGGER.

The `table` object supports the following:

* `database_name` - (Required, String) Database name.
* `privileges` - (Required, Set) Table privilege.available values for Privileges:SELECT,INSERT,UPDATE,DELETE,CREATE,DROP,REFERENCES,INDEX,ALTER,CREATE VIEW,SHOW VIEW,TRIGGER.
* `table_name` - (Required, String) Table name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



