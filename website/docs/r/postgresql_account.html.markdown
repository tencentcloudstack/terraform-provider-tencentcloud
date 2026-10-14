---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_account"
sidebar_current: "docs-tencentcloud-resource-postgresql_account"
description: |-
  Provides a resource to create a postgresql account
---

# tencentcloud_postgresql_account

Provides a resource to create a postgresql account

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
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance ID in the format of postgres-4wdeb0zv.
* `password` - (Required, String) Password, which can contain 8-32 letters, digits, and symbols (()`~!@#$%^&amp;amp;amp;*-+=_|{}[]:;&amp;amp;#39;&amp;amp;lt;&amp;amp;gt;,.?/); can&amp;amp;#39;t start with slash /.
* `type` - (Required, String, ForceNew) The type of user. Valid values: 1. normal: regular user; 2. tencentDBSuper: user with the pg_tencentdb_superuser role.
* `user_name` - (Required, String, ForceNew) Instance username, which can contain 1-16 letters, digits, and underscore (_); can&amp;amp;#39;t be postgres; can&amp;amp;#39;t start with numbers, pg_, and tencentdb_.
* `lock_status` - (Optional, Bool) whether lock account. true: locked; false: unlock.
* `remark` - (Optional, String) Remarks correspond to user `UserName`, which can contain 0-60 letters, digits, symbols (-_), and Chinese characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAccount](https://www.tencentcloud.com/document/product/409/78820) - This API is used to create a data account. The returned Oid is the unique identifier of the account, which matches the oid recorded in the database system table pg_roles.
* Create: [LockAccount](https://www.tencentcloud.com/document/product/409/78817) - This API is used to lock a database account. After being locked, the current connection will disconnect and unable to establish a new connection.
* Read: [DescribeAccounts](https://www.tencentcloud.com/document/product/409/18109) - This API is used to query the list of the database accounts for an instance.
* Update: [LockAccount](https://www.tencentcloud.com/document/product/409/78817) - This API is used to lock a database account. After being locked, the current connection will disconnect and unable to establish a new connection.
* Update: [ModifyAccountRemark](https://www.tencentcloud.com/document/product/409/18108) - This API is used to modify account remarks.
* Update: [ResetAccountPassword](https://www.tencentcloud.com/document/product/409/18107) - This API is used to reset the account password of an instance.
* Update: [UnlockAccount](https://www.tencentcloud.com/document/product/409/78815) - This API is used to unlock a database account. The account can log in to the database after being unlocked.
* Delete: [DeleteAccount](https://www.tencentcloud.com/document/product/409/78819) - This API is used to delete a database account. Oid and UserName must be input simultaneously to avoid accidental deletion. Note: This API is reentrant.

## Import

postgres account can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_account.example postgres-3hk6b6tj#tf_example
```

