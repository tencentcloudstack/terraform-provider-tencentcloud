---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_backup_config"
sidebar_current: "docs-tencentcloud-resource-cynosdb_backup_config"
description: |-
  Provides a resource to create a cynosdb backup_config
---

# tencentcloud_cynosdb_backup_config

Provides a resource to create a cynosdb backup_config

## Example Usage

### Enable logical backup configuration and cross-region logical backup

```hcl
resource "tencentcloud_cynosdb_backup_config" "foo" {
  backup_time_beg  = 7200
  backup_time_end  = 21600
  cluster_id       = "cynosdbmysql-bws8h88b"
  reserve_duration = 604800

  logic_backup_config {
    logic_backup_enable        = "ON"
    logic_backup_time_beg      = 7200
    logic_backup_time_end      = 21600
    logic_cross_regions        = ["ap-shanghai"]
    logic_cross_regions_enable = "ON"
    logic_reserve_duration     = 259200
  }
}
```

### Disable logical backup configuration

```hcl
resource "tencentcloud_cynosdb_backup_config" "foo" {
  backup_time_beg  = 7200
  backup_time_end  = 21600
  cluster_id       = "cynosdbmysql-bws8h88b"
  reserve_duration = 604800

  logic_backup_config {
    logic_backup_enable = "OFF"
  }
}
```

## Argument Reference

The following arguments are supported:

* `backup_time_beg` - (Required, Int) Full backup start time. Value range: [0-24*3600]. For example, 0:00 AM, 1:00 AM, and 2:00 AM are represented by 0, 3600, and 7200, respectively.
* `backup_time_end` - (Required, Int) Full backup end time. Value range: [0-24*3600]. For example, 0:00 AM, 1:00 AM, and 2:00 AM are represented by 0, 3600, and 7200, respectively.
* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `reserve_duration` - (Required, Int) Backup retention period in seconds. Backups will be cleared after this period elapses. 7 days is represented by 3600*24*7 = 604800. Maximum value: 158112000.
* `logic_backup_config` - (Optional, List) Logical backup configuration. Do not set this field if it is not enabled. Example value: [{"LogicBackupEnable": "ON","LogicBackupTimeBeg": "2023-04-24 15:06:04","LogicBackupTimeEnd": "2024-04-24 15:06:04","LogicReserveDuration": "60","LogicCrossRegionsEnable": "ON","LogicCrossRegions": ["ap-guangzhou"]}].

The `logic_backup_config` object supports the following:

* `logic_backup_enable` - (Optional, String) Whether to enable automatic logical backup. Value: `ON`, `OFF`.
* `logic_backup_time_beg` - (Optional, Int) Automatic logical backup start time. When `logic_backup_enable` is `OFF`, it must be `0` or not entered. Example value: 2.
* `logic_backup_time_end` - (Optional, Int) Automatic logical backup end time. When `logic_backup_enable` is `OFF`, it must be `0` or not entered. Example value: 6.
* `logic_cross_regions_enable` - (Optional, String) Whether to enable cross-region logical backup. Cannot be input when `logic_backup_enable` is `OFF`. When `logic_backup_enable` is `ON`, `logic_cross_regions_enable` setting `ON` will take effect. Value: `ON`, `OFF`.
* `logic_cross_regions` - (Optional, Set) Logical backup across regions. Example value: ["ap-guangzhou"]. When `logic_backup_enable` is `OFF`, it must be `[]` or not entered.
* `logic_reserve_duration` - (Optional, Int) Automatic logical backup retention period. When `logic_backup_enable` is `OFF`, it must be `0` or not entered. Value range: [259200,158112000]. `logic_backup_enable` is `OFF`, `logic_reserve_duration` cannot be set when creating.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `backup_freq` - Backup frequency. It is an array of 7 elements corresponding to Monday through Sunday. full: full backup; increment: incremental backup. This parameter cannot be modified currently and doesn't need to be entered.
* `backup_type` - Backup mode. logic: logic backup; snapshot: snapshot backup. This parameter cannot be modified currently and doesn't need to be entered.


## Import

cynosdb backup_config can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_backup_config.foo cynosdbmysql-bws8h88b
```

