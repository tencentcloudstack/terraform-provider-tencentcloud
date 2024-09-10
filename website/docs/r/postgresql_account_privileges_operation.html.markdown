---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_account_privileges_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_account_privileges_operation"
description: |-
  Provides a resource to create postgresql account privileges
---

# tencentcloud_postgresql_account_privileges_operation

Provides a resource to create postgresql account privileges

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

# create postgresql
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

# create account
resource "tencentcloud_postgresql_account" "example" {
  db_instance_id = tencentcloud_postgresql_instance.example.id
  user_name      = "tf_example"
  password       = "Password@123"
  type           = "normal"
  remark         = "remark"
  lock_status    = false
}

# create account privileges
resource "tencentcloud_postgresql_account_privileges_operation" "example" {
  db_instance_id = tencentcloud_postgresql_instance.example.id
  user_name      = tencentcloud_postgresql_account.example.user_name
  modify_privilege_set {
    database_privilege {
      object {
        object_name = "postgres"
        object_type = "database"
      }

      privilege_set = ["CONNECT", "TEMPORARY", "CREATE"]
    }

    modify_type = "grantObject"
    is_cascade  = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance ID in the format of postgres-4wdeb0zv.
* `modify_privilege_set` - (Required, List) Privileges to modify. Batch modification supported, up to 50 entries at a time.
* `user_name` - (Required, String, ForceNew) Instance username.

The `database_privilege` object of `modify_privilege_set` supports the following:

* `object` - (Optional, List) Database object.If ObjectType is database, DatabaseName/SchemaName/TableName can be null.If ObjectType is schema, SchemaName/TableName can be null.If ObjectType is table, TableName can be null.If ObjectType is column, DatabaseName/SchemaName/TableName can&amp;#39;t be null.In all other cases, DatabaseName/SchemaName/TableName can be null. Note: This field may return null, indicating that no valid value can be obtained.
* `privilege_set` - (Optional, Set) Privileges the specific account has on database object. Note: This field may return null, indicating that no valid value can be obtained.

The `modify_privilege_set` object supports the following:

* `database_privilege` - (Optional, List) Database objects and the user permissions on these objects. Note: This field may return null, indicating that no valid value can be obtained.
* `is_cascade` - (Optional, Bool) Required only when ModifyType is revokeObject. When the parameter is true, revoking permissions will cascade. The default value is false.
* `modify_type` - (Optional, String) Supported modification method: grantObject, revokeObject, alterRole. grantObject represents granting permissions on object, revokeObject represents revoking permissions on object, and alterRole represents modifying the account type.

The `object` object of `database_privilege` supports the following:

* `object_name` - (Required, String) Database object Name. Note: This field may return null, indicating that no valid value can be obtained.
* `object_type` - (Required, String) Supported database object types: account, database, schema, sequence, procedure, type, function, table, view, matview, column. Note: This field may return null, indicating that no valid value can be obtained.
* `database_name` - (Optional, String) Database name to which the database object belongs. This parameter is mandatory when ObjectType is not database. Note: This field may return null, indicating that no valid value can be obtained.
* `schema_name` - (Optional, String) Schema name to which the database object belongs. This parameter is mandatory when ObjectType is not database or schema. Note: This field may return null, indicating that no valid value can be obtained.
* `table_name` - (Optional, String) Table name to which the database object belongs. This parameter is mandatory when ObjectType is column. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



