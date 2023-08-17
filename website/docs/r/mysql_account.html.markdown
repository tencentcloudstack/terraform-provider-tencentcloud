---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_account"
sidebar_current: "docs-tencentcloud-resource-mysql_account"
description: |-
  Provides a MySQL account resource for database management. A MySQL instance supports multiple database account.
---

# tencentcloud_mysql_account

Provides a MySQL account resource for database management. A MySQL instance supports multiple database account.

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
```

## Argument Reference

The following arguments are supported:

* `mysql_id` - (Required, String, ForceNew) Instance ID to which the account belongs.
* `name` - (Required, String, ForceNew) Account name.
* `password` - (Required, String) Operation password.
* `description` - (Optional, String) Database description.
* `host` - (Optional, String) Account host, default is `%`.
* `max_user_connections` - (Optional, Int) The maximum number of available connections for a new account, the default value is 10240, and the maximum value that can be set is 10240.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql account can be imported using the mysqlId#accountName, e.g.

```
terraform import tencentcloud_mysql_account.default cdb-gqg6j82x#tf_account
```

