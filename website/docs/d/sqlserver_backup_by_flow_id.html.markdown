---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_backup_by_flow_id"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_backup_by_flow_id"
description: |-
  Use this data source to query detailed information of sqlserver datasource_backup_by_flow_id
---

# tencentcloud_sqlserver_backup_by_flow_id

Use this data source to query detailed information of sqlserver datasource_backup_by_flow_id

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

data "tencentcloud_sqlserver_backup_by_flow_id" "example" {
  instance_id = tencentcloud_sqlserver_general_backup.example.instance_id
  flow_id     = tencentcloud_sqlserver_general_backup.example.flow_id
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
  instance_id = tencentcloud_sqlserver_db.example.id
  backup_name = "tf_example_backup"
  strategy    = 0
}
```

## Argument Reference

The following arguments are supported:

* `flow_id` - (Required, String) Create a backup process ID, which can be obtained through the [CreateBackup](https://cloud.tencent.com/document/product/238/19946) interface.
* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backup_name` - Backup task name, customizable.
* `backup_way` - Backup method, 0-scheduled backup; 1-manual temporary backup; instance status is 0-creating, this field is the default value 0, meaningless.
* `dbs` - For the DB list, only the library name contained in the first record is returned for a single-database backup file; for a single-database backup file, the library names of all records need to be obtained through the DescribeBackupFiles interface.
* `end_time` - backup end time.
* `external_addr` - External network download address, for a single database backup file, only the external network download address of the first record is returned; single database backup files need to obtain the download addresses of all records through the DescribeBackupFiles interface.
* `file_name` - File name. For a single-database backup file, only the file name of the first record is returned; for a single-database backup file, the file names of all records need to be obtained through the DescribeBackupFiles interface.
* `group_id` - Aggregate Id, this value is not returned for packaged backup files. Use this value to call the DescribeBackupFiles interface to obtain the detailed information of a single database backup file.
* `internal_addr` - Intranet download address, for a single database backup file, only the intranet download address of the first record is returned; single database backup files need to obtain the download addresses of all records through the DescribeBackupFiles interface.
* `start_time` - backup start time.
* `status` - Backup file status, 0-creating; 1-success; 2-failure.
* `strategy` - Backup strategy, 0-instance backup; 1-multi-database backup; when the instance status is 0-creating, this field is the default value 0, meaningless.


