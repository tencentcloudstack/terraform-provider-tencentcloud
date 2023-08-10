---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_rollback"
sidebar_current: "docs-tencentcloud-resource-mysql_rollback"
description: |-
  Provides a resource to create a mysql rollback
---

# tencentcloud_mysql_rollback

Provides a resource to create a mysql rollback

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

data "tencentcloud_mysql_rollback_range_time" "example" {
  instance_ids = [tencentcloud_mysql_instance.example.id]
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
  slave_deploy_mode = 1
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  first_slave_zone  = data.tencentcloud_availability_zones_by_product.zones.zones.1.name
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

resource "tencentcloud_mysql_rollback" "example" {
  instance_id   = tencentcloud_mysql_instance.example.id
  strategy      = "full"
  rollback_time = data.tencentcloud_mysql_rollback_range_time.example.item.0.times.0.start
  databases {
    database_name     = "tf_db_bak"
    new_database_name = "tf_db_bak_new"
  }
  tables {
    database = "tf_db_bak1"
    table {
      table_name     = "tf_table"
      new_table_name = "tf_table_new"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Cloud database instance ID.
* `rollback_time` - (Required, String, ForceNew) Database rollback time, the time format is: yyyy-mm-dd hh:mm:ss.
* `strategy` - (Required, String, ForceNew) Rollback strategy. Available values are: table, db, full; the default value is full. table- Extremely fast rollback mode, only import the backup and binlog of the selected table level, if there is a cross-table operation, and the associated table is not selected at the same time, the rollback will fail. In this mode, the parameter Databases must be empty; db- Quick mode, only import the backup and binlog of the selected library level, if there is a cross-database operation, and the associated library is not selected at the same time, the rollback will fail; full- normal rollback mode, the backup and binlog of the entire instance will be imported, at a slower rate.
* `databases` - (Optional, List, ForceNew) The database information to be archived, indicating that the entire database is archived.
* `tables` - (Optional, List, ForceNew) The database table information to be rolled back, indicating that the file is rolled back by table.

The `databases` object supports the following:

* `database_name` - (Required, String) The original database name before rollback.
* `new_database_name` - (Required, String) The new database name after rollback.

The `table` object supports the following:

* `new_table_name` - (Required, String) New database table name after rollback.
* `table_name` - (Required, String) The original database table name before rollback.

The `tables` object supports the following:

* `database` - (Required, String) Database name.
* `table` - (Required, List) Database table details.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



