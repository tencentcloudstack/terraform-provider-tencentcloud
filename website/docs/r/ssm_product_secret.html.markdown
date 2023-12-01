---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_product_secret"
sidebar_current: "docs-tencentcloud-resource-ssm_product_secret"
description: |-
  Provides a resource to create a ssm product_secret
---

# tencentcloud_ssm_product_secret

Provides a resource to create a ssm product_secret

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 0
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  slave_sync_mode   = 1
  instance_name     = "tf-example"
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [tencentcloud_security_group.security_group.id]

  tags = {
    createBy = "terraform"
  }

  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_kms_key" "example" {
  alias                = "tf-example-kms-key"
  description          = "example of kms key"
  key_rotation_enabled = false
  is_enabled           = true

  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_ssm_product_secret" "example" {
  secret_name      = "tf-example"
  user_name_prefix = "prefix"
  product_name     = "Mysql"
  instance_id      = tencentcloud_mysql_instance.example.id
  domains          = ["10.0.0.0"]
  privileges_list {
    privilege_name = "GlobalPrivileges"
    privileges     = ["ALTER ROUTINE"]
  }
  description         = "for ssm product test"
  kms_key_id          = tencentcloud_kms_key.example.id
  status              = "Enabled"
  enable_rotation     = true
  rotation_begin_time = "2023-08-05 20:54:33"
  rotation_frequency  = 30

  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domains` - (Required, Set: [`String`]) Domain name of the account in the form of IP. You can enter `%`.
* `instance_id` - (Required, String) Tencent Cloud service instance ID.
* `privileges_list` - (Required, List) List of permissions that need to be granted when the credential is bound to a Tencent Cloud service.
* `product_name` - (Required, String) Name of the Tencent Cloud service bound to the credential, such as `Mysql`, `Tdsql-mysql`. you can use dataSource `tencentcloud_ssm_products` to query supported products.
* `secret_name` - (Required, String, ForceNew) Credential name, which must be unique in the same region. It can contain 128 bytes of letters, digits, hyphens, and underscores and must begin with a letter or digit.
* `user_name_prefix` - (Required, String) Prefix of the user account name, which is specified by you and can contain up to 8 characters.Supported character sets include:Digits: [0, 9].Lowercase letters: [a, z].Uppercase letters: [A, Z].Special symbols: underscore.The prefix must begin with a letter.
* `description` - (Optional, String) Description, which is used to describe the purpose in detail and can contain up to 2,048 bytes.
* `enable_rotation` - (Optional, Bool) Specifies whether to enable rotation, when secret status is `Disabled`, rotation will be disabled. `True` - enable, `False` - do not enable. If this parameter is not specified, `False` will be used by default.
* `kms_key_id` - (Optional, String) Specifies the KMS CMK that encrypts the credential. If this parameter is left empty, the CMK created by Secrets Manager by default will be used for encryption.You can also specify a custom KMS CMK created in the same region for encryption.
* `rotation_begin_time` - (Optional, String) User-Defined rotation start time in the format of 2006-01-02 15:04:05.When `EnableRotation` is `True`, this parameter is required.
* `rotation_frequency` - (Optional, Int) Rotation frequency in days. Default value: 1 day.
* `status` - (Optional, String) Enable or Disable Secret. Valid values is `Enabled` or `Disabled`. Default is `Enabled`.
* `tags` - (Optional, Map) Tags of secret.

The `privileges_list` object supports the following:

* `privilege_name` - (Required, String) Permission name. Valid values: `GlobalPrivileges`, `DatabasePrivileges`, `TablePrivileges`, `ColumnPrivileges`. When the permission is `DatabasePrivileges`, the database name must be specified by the `Database` parameter; When the permission is `TablePrivileges`, the database name and the table name in the database must be specified by the `Database` and `TableName` parameters; When the permission is `ColumnPrivileges`, the database name, table name in the database, and column name in the table must be specified by the `Database`, `TableName`, and `ColumnName` parameters.
* `privileges` - (Required, Set) Permission list. For the `Mysql` service, optional permission values are: 1. Valid values of `GlobalPrivileges`: SELECT,INSERT,UPDATE,DELETE,CREATE, PROCESS, DROP,REFERENCES,INDEX,ALTER,SHOW DATABASES,CREATE TEMPORARY TABLES,LOCK TABLES,EXECUTE,CREATE VIEW,SHOW VIEW,CREATE ROUTINE,ALTER ROUTINE,EVENT,TRIGGER. Note: if this parameter is not passed in, it means to clear the permission. 2. Valid values of `DatabasePrivileges`: SELECT,INSERT,UPDATE,DELETE,CREATE, DROP,REFERENCES,INDEX,ALTER,CREATE TEMPORARY TABLES,LOCK TABLES,EXECUTE,CREATE VIEW,SHOW VIEW,CREATE ROUTINE,ALTER ROUTINE,EVENT,TRIGGER. Note: if this parameter is not passed in, it means to clear the permission. 3. Valid values of `TablePrivileges`: SELECT,INSERT,UPDATE,DELETE,CREATE, DROP,REFERENCES,INDEX,ALTER,CREATE VIEW,SHOW VIEW, TRIGGER. Note: if this parameter is not passed in, it means to clear the permission. 4. Valid values of `ColumnPrivileges`: SELECT,INSERT,UPDATE,REFERENCES.Note: if this parameter is not passed in, it means to clear the permission.
* `column_name` - (Optional, String) This value takes effect only when `PrivilegeName` is `ColumnPrivileges`, and the following parameters are required in this case:Database: explicitly indicate the database instance.TableName: explicitly indicate the table.
* `database` - (Optional, String) This value takes effect only when `PrivilegeName` is `DatabasePrivileges`.
* `table_name` - (Optional, String) This value takes effect only when `PrivilegeName` is `TablePrivileges`, and the `Database` parameter is required in this case to explicitly indicate the database instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Credential creation time in UNIX timestamp format.
* `secret_type` - `0`: user-defined secret. `1`: Tencent Cloud services secret. `2`: SSH key secret. `3`: Tencent Cloud API key secret. Note: this field may return `null`, indicating that no valid values can be obtained.


