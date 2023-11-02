---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_backup_policy"
sidebar_current: "docs-tencentcloud-resource-mysql_backup_policy"
description: |-
  Provides a mysql policy resource to create a backup policy.
---

# tencentcloud_mysql_backup_policy

Provides a mysql policy resource to create a backup policy.

~> **NOTE:** This attribute `backup_model` only support 'physical' in Terraform TencentCloud provider version 1.16.2

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

resource "tencentcloud_mysql_backup_policy" "example" {
  mysql_id              = tencentcloud_mysql_instance.example.id
  retention_period      = 7
  backup_model          = "physical"
  backup_time           = "22:00-02:00"
  binlog_period         = 32
  enable_binlog_standby = "off"
  binlog_standby_days   = 31
}
```

## Argument Reference

The following arguments are supported:

* `mysql_id` - (Required, String, ForceNew) Instance ID to which policies will be applied.
* `backup_model` - (Optional, String) Backup method. Supported values include: `physical` - physical backup.
* `backup_time` - (Optional, String) Instance backup time, in the format of 'HH:mm-HH:mm'. Time setting interval is four hours. Default to `02:00-06:00`. The following value can be supported: `02:00-06:00`, `06:00-10:00`, `10:00-14:00`, `14:00-18:00`, `18:00-22:00`, and `22:00-02:00`.
* `binlog_period` - (Optional, Int) Binlog retention time, in days. The minimum value is 7 days and the maximum value is 1830 days. This value cannot be set greater than the backup file retention time.
* `binlog_standby_days` - (Optional, Int) The standard starting number of days for log backup storage. The log backup will be converted when it reaches the standard starting number of days for storage. The minimum is 30 days and must not be greater than the number of days for log backup retention.
* `enable_binlog_standby` - (Optional, String) Whether to enable the log backup standard storage policy, `off` - close, `on` - open, the default is off.
* `retention_period` - (Optional, Int) The retention time of backup files, in days. The minimum value is 7 days and the maximum value is 1830 days. And default value is `7`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



