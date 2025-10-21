---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_backup_strategy"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_backup_strategy"
description: |-
  Provides a resource to create a sqlserver config_backup_strategy
---

# tencentcloud_sqlserver_config_backup_strategy

Provides a resource to create a sqlserver config_backup_strategy

## Example Usage

### Daily backup

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
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

resource "tencentcloud_sqlserver_config_backup_strategy" "example" {
  instance_id              = tencentcloud_sqlserver_basic_instance.example.id
  backup_type              = "daily"
  backup_time              = 0
  backup_day               = 1
  backup_model             = "master_no_pkg"
  backup_cycle             = [1]
  backup_save_days         = 7
  regular_backup_enable    = "disable"
  regular_backup_save_days = 90
  regular_backup_strategy  = "months"
  regular_backup_counts    = 1
}
```

### Weekly backup

```hcl
resource "tencentcloud_sqlserver_config_backup_strategy" "example" {
  instance_id              = tencentcloud_sqlserver_basic_instance.example.id
  backup_type              = "weekly"
  backup_time              = 0
  backup_model             = "master_no_pkg"
  backup_cycle             = [1, 3, 5]
  backup_save_days         = 7
  regular_backup_enable    = "disable"
  regular_backup_save_days = 90
  regular_backup_strategy  = "months"
  regular_backup_counts    = 1
}
```

### Regular backup

```hcl
resource "tencentcloud_sqlserver_config_backup_strategy" "example" {
  instance_id               = tencentcloud_sqlserver_basic_instance.example.id
  backup_time               = 0
  backup_model              = "master_no_pkg"
  backup_cycle              = [1, 3]
  backup_save_days          = 7
  regular_backup_enable     = "enable"
  regular_backup_save_days  = 120
  regular_backup_strategy   = "months"
  regular_backup_counts     = 1
  regular_backup_start_time = "%s"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `backup_cycle` - (Optional, Set: [`Int`]) The days of the week on which backup will be performed when `BackupType` is weekly. If data backup retention period is less than 7 days, the values will be 1-7, indicating that backup will be performed everyday by default; if data backup retention period is greater than or equal to 7 days, the values will be at least any two days, indicating that backup will be performed at least twice in a week by default.
* `backup_day` - (Optional, Int) Backup interval in days when the BackupType is daily. The current value can only be 1.
* `backup_model` - (Optional, String) Backup mode. Valid values: master_pkg (archive the backup files of the primary node), master_no_pkg (do not archive the backup files of the primary node), slave_pkg (archive the backup files of the replica node), slave_no_pkg (do not archive the backup files of the replica node). Backup files of the replica node are supported only when Always On disaster recovery is enabled.
* `backup_save_days` - (Optional, Int) Data (log) backup retention period. Value range: 3-1830 days, default value: 7 days.
* `backup_time` - (Optional, Int) Backup time. Value range: an integer from 0 to 23.
* `backup_type` - (Optional, String) Backup type. Valid values: weekly (when length(BackupDay) <=7 && length(BackupDay) >=2), daily (when length(BackupDay)=1). Default value: daily.
* `regular_backup_counts` - (Optional, Int) The number of retained archive backups. Default value: 1.
* `regular_backup_enable` - (Optional, String) Archive backup status. Valid values: enable (enabled); disable (disabled). Default value: disable.
* `regular_backup_save_days` - (Optional, Int) Archive backup retention days. Value range: 90-3650 days. Default value: 365 days.
* `regular_backup_start_time` - (Optional, String) Archive backup start date in YYYY-MM-DD format, which is the current time by default.
* `regular_backup_strategy` - (Optional, String) Archive backup policy. Valid values: years (yearly); quarters (quarterly); months(monthly); Default value: `months`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver config_backup_strategy can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_backup_strategy.example mssql-si2823jyl
```

