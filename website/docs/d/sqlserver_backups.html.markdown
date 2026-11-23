---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_backups"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_backups"
description: |-
  Use this data source to query the list of SQL Server backups.
---

# tencentcloud_sqlserver_backups

Use this data source to query the list of SQL Server backups.

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

data "tencentcloud_sqlserver_backups" "example" {
  instance_id = tencentcloud_sqlserver_general_backup.example.instance_id
  start_time  = "2023-08-01 00:00:00"
  end_time    = "2023-08-07 00:00:00"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_db"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_general_backup" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
  backup_name = "tf_example_backup"
  strategy    = 0
}
```

### Filter by backup name

```hcl
data "tencentcloud_sqlserver_backups" "example" {
  instance_id = tencentcloud_sqlserver_general_backup.example.instance_id
  start_time  = "2023-08-01 00:00:00"
  end_time    = "2023-08-07 00:00:00"
  backup_name = "tf-example-backup"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time of the instance list, like yyyy-MM-dd HH:mm:ss.
* `instance_id` - (Required, String) Instance ID.
* `start_time` - (Required, String) Start time of the instance list, like yyyy-MM-dd HH:mm:ss.
* `backup_name` - (Optional, String) Filter by backup name, do not filter if left blank.
* `result_output_file` - (Optional, String) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of SQL Server backup. Each element contains the following attributes:
  * `db_list` - Database name list of the backup.
  * `end_time` - End time of the backup.
  * `file_name` - File name of the backup.
  * `id` - ID of the backup.
  * `instance_id` - Instance ID.
  * `internet_url` - URL for downloads externally.
  * `intranet_url` - URL for downloads internally.
  * `size` - The size of backup file. Unit is KB.
  * `start_time` - Start time of the backup.
  * `status` - Status of the backup. `1` for creating, `2` for successfully created, 3 for failed.
  * `strategy` - Strategy of the backup. `0` for instance backup, `1` for multi-databases backup.
  * `trigger_model` - The way to trigger backup. `0` for timed trigger, `1` for manual trigger.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeBackups](https://www.tencentcloud.com/document/api/238/32127)

