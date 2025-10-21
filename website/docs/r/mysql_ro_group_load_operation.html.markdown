---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_ro_group_load_operation"
sidebar_current: "docs-tencentcloud-resource-mysql_ro_group_load_operation"
description: |-
  Provides a resource to create a mysql ro_group_load_operation
---

# tencentcloud_mysql_ro_group_load_operation

Provides a resource to create a mysql ro_group_load_operation

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

data "tencentcloud_mysql_instance" "example" {
  mysql_id = tencentcloud_mysql_instance.example.id

  depends_on = [tencentcloud_mysql_readonly_instance.example]
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.1.name
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
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.1.name
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

resource "tencentcloud_mysql_readonly_instance" "example" {
  master_instance_id = tencentcloud_mysql_instance.example.id
  instance_name      = "tf-mysql"
  mem_size           = 2000
  volume_size        = 200
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  intranet_port      = 3306
  security_groups    = [tencentcloud_security_group.security_group.id]

  tags = {
    createBy = "terraform"
  }
}

resource "tencentcloud_mysql_ro_group_load_operation" "ro_group_load_operation" {
  ro_group_id = data.tencentcloud_mysql_instance.example.instance_list.0.ro_groups.0.group_id
}
```

## Argument Reference

The following arguments are supported:

* `ro_group_id` - (Required, String, ForceNew) The ID of the RO group, in the format: cdbrg-c1nl9rpv.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



