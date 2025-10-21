---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_remote_backup_config"
sidebar_current: "docs-tencentcloud-resource-mysql_remote_backup_config"
description: |-
  Provides a resource to create a mysql remote_backup_config
---

# tencentcloud_mysql_remote_backup_config

Provides a resource to create a mysql remote_backup_config

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

resource "tencentcloud_mysql_remote_backup_config" "example" {
  instance_id        = tencentcloud_mysql_instance.example.id
  remote_backup_save = "on"
  remote_binlog_save = "on"
  remote_region      = ["ap-shanghai"]
  expire_days        = 7
}
```

## Argument Reference

The following arguments are supported:

* `expire_days` - (Required, Int) Remote backup retention time, in days.
* `instance_id` - (Required, String) Instance ID, in the format: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.
* `remote_backup_save` - (Required, String) Remote data backup switch, off - disable remote backup, on - enable remote backup.
* `remote_binlog_save` - (Required, String) Off-site log backup switch, off - off off-site backup, on-on off-site backup, only when the parameter RemoteBackupSave is on, the RemoteBinlogSave parameter can be set to on.
* `remote_region` - (Required, Set: [`String`]) User settings off-site backup region list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql remote_backup_config can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_remote_backup_config.remote_backup_config remote_backup_config_id
```

