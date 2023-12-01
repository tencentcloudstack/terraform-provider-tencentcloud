---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_backup_encryption_status"
sidebar_current: "docs-tencentcloud-resource-mysql_backup_encryption_status"
description: |-
  Provides a resource to create a mysql backup_encryption_status
---

# tencentcloud_mysql_backup_encryption_status

Provides a resource to create a mysql backup_encryption_status

## Example Usage

### Enable encryption

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

resource "tencentcloud_mysql_backup_encryption_status" "example" {
  instance_id       = tencentcloud_mysql_instance.example.id
  encryption_status = "on"
}
```

### Disable encryption

```hcl
resource "tencentcloud_mysql_backup_encryption_status" "example" {
  instance_id       = tencentcloud_mysql_instance.example.id
  encryption_status = "off"
}
```

## Argument Reference

The following arguments are supported:

* `encryption_status` - (Required, String) Whether physical backup encryption is enabled for the instance. Possible values are `on`, `off`.
* `instance_id` - (Required, String) Instance ID, in the format: cdb-XXXX. Same instance ID as displayed in the ApsaraDB for Console page.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql backup_encryption_status can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_backup_encryption_status.backup_encryption_status backup_encryption_status_id
```

