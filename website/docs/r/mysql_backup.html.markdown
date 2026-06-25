---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_backup"
sidebar_current: "docs-tencentcloud-resource-mysql_backup"
description: |-
  Provides a resource to create a CDB mysql backup
---

# tencentcloud_mysql_backup

Provides a resource to create a CDB mysql backup

## Example Usage

### Create a physical full backup

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

resource "tencentcloud_mysql_backup" "example" {
  instance_id        = tencentcloud_mysql_instance.example.id
  backup_method      = "physical"
  manual_backup_name = "tf-example-backup"
  encryption_flag    = "off"
}
```

### Create a logical backup with specific database and table

```hcl
resource "tencentcloud_mysql_backup" "logical" {
  instance_id        = tencentcloud_mysql_instance.example.id
  backup_method      = "logical"
  manual_backup_name = "tf-logical-backup"

  backup_db_table_list {
    database = "db1"
    table    = "tb1"
  }

  backup_db_table_list {
    database = "db2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `backup_method` - (Required, String, ForceNew) Target backup method. Supported values include: `logical` - logical cold backup, `physical` - physical cold backup, `snapshot` - snapshot backup. Basic edition instances only support snapshot backup.
* `instance_id` - (Required, String, ForceNew) Instance ID, such as `cdb-c1nl9rpv`. It is identical to the instance ID displayed in the database console page.
* `backup_db_table_list` - (Optional, List, ForceNew) List of databases and tables to backup. Only valid when `backup_method` is `logical`. The specified databases and tables must exist, otherwise backup may fail.
* `encryption_flag` - (Optional, String, ForceNew) Whether to encrypt physical backup. Supported values include: `on` - yes, `off` - no. Only valid when `backup_method` is `physical`. If not specified, the instance's default backup encryption policy is used.
* `manual_backup_name` - (Optional, String, ForceNew) Manual backup alias. Maximum length is 60 characters.

The `backup_db_table_list` object supports the following:

* `database` - (Required, String) Database name.
* `table` - (Optional, String) Table name. If specified, backup this table in the database. If not specified, backup the entire database.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `backup_id` - ID of the backup task.


## Import

mysql backup can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_backup.foo backupId#instanceId
```

