---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_audit_log_file"
sidebar_current: "docs-tencentcloud-resource-mysql_audit_log_file"
description: |-
  Provides a resource to create a mysql audit_log_file
---

# tencentcloud_mysql_audit_log_file

Provides a resource to create a mysql audit_log_file

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

resource "tencentcloud_mysql_audit_log_file" "example" {
  instance_id = tencentcloud_mysql_instance.example.id
  start_time  = "2023-07-01 00:00:00"
  end_time    = "2023-10-01 00:00:00"
  order       = "ASC"
  order_by    = "timestamp"
}
```

### Add filter

```hcl
resource "tencentcloud_mysql_audit_log_file" "example" {
  instance_id = tencentcloud_mysql_instance.example.id
  start_time  = "2023-07-01 00:00:00"
  end_time    = "2023-10-01 00:00:00"
  order       = "ASC"
  order_by    = "timestamp"

  filter {
    host = ["30.50.207.46"]
    user = ["keep_dbbrain"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String, ForceNew) end time.
* `instance_id` - (Required, String, ForceNew) The ID of instance.
* `start_time` - (Required, String, ForceNew) start time.
* `filter` - (Optional, List, ForceNew) Filter condition. Logs can be filtered according to the filter conditions set.
* `order_by` - (Optional, String, ForceNew) Sort field. supported values include:`timestamp` - timestamp; `affectRows` - affected rows; `execTime` - execution time.
* `order` - (Optional, String, ForceNew) Sort by. supported values are: `ASC`- ascending order, `DESC`- descending order.

The `filter` object supports the following:

* `affect_rows` - (Optional, Int) Affects the number of rows. Indicates to filter audit logs whose number of affected rows is greater than this value.
* `db_name` - (Optional, Set) Database name.
* `exec_time` - (Optional, Int) Execution time. The unit is: ms. Indicates to filter audit logs whose execution time is greater than this value.
* `host` - (Optional, Set) Client address.
* `policy_name` - (Optional, Set) The name of policy.
* `sql_type` - (Optional, String) SQL type. Currently supported: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, SET, REPLACE, EXECUTE.
* `sql_types` - (Optional, Set) SQL type. Supports simultaneous query of multiple types. Currently supported: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, ALTER, SET, REPLACE, EXECUTE.
* `sql` - (Optional, String) SQL statement. support fuzzy matching.
* `sqls` - (Optional, Set) SQL statement. Support passing multiple sql statements.
* `table_name` - (Optional, Set) Table name.
* `user` - (Optional, Set) User name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `download_url` - download url.
* `file_size` - size of file(KB).


