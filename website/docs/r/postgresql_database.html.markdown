---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_database"
sidebar_current: "docs-tencentcloud-resource-postgresql_database"
description: |-
  Provides a resource to manage a TencentDB for PostgreSQL database within a DB instance.
---

# tencentcloud_postgresql_database

Provides a resource to manage a TencentDB for PostgreSQL database within a DB instance.

-> **Note:** The `database_owner` must be an existing account of the PostgreSQL instance. You can use resource `tencentcloud_postgresql_account` to create the account before creating the database.

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql instance
resource "tencentcloud_postgresql_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  db_major_version  = "10"
  engine_version    = "10.23"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  cpu               = 1
  memory            = 2
  storage           = 10

  tags = {
    test = "tf"
  }
}

# create postgresql account
resource "tencentcloud_postgresql_account" "example" {
  db_instance_id = tencentcloud_postgresql_instance.example.id
  user_name      = "tf_example"
  password       = "Password@123"
  type           = "normal"
  remark         = "remark"
  lock_status    = false
}

# create postgresql database
resource "tencentcloud_postgresql_database" "example" {
  db_instance_id = tencentcloud_postgresql_instance.example.id
  database_name  = "test_db"
  database_owner = tencentcloud_postgresql_account.example.user_name
  encoding       = "UTF8"
  collate        = "C"
  ctype          = "C"
}
```

## Argument Reference

The following arguments are supported:

* `database_name` - (Required, String, ForceNew) Database name.
* `database_owner` - (Required, String) Database owner account.
* `db_instance_id` - (Required, String, ForceNew) DB instance ID, such as postgres-6fego161.
* `collate` - (Optional, String, ForceNew) Database collation rule.
* `ctype` - (Optional, String, ForceNew) Database character classification.
* `encoding` - (Optional, String, ForceNew) Database character encoding, such as UTF8, LATIN1, LATIN2, WIN1250, WIN1251, WIN1252, KOI8R, EUC_JP, EUC_KR. Default: UTF8.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

PostgreSQL database can be imported using the composite ID `db_instance_id#database_name`, e.g.

```
terraform import tencentcloud_postgresql_database.example postgres-6fego161#test_db
```

