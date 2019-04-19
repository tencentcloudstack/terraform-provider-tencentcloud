---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_backup_policy"
sidebar_current: "docs-tencentcloud-tencentcloud_mysql_backup_policy"
description: |-
Provides a mysql policy resource to create a backup policy.

---
#tencentcloud_mysql_backup_policy

Provides a mysql policy resource to create a backup policy.

##Example Usage

```
resource " tencentcloud_mysql_backup_policy " "default" {
  mysql_id = "cdb-dnqksd9f"
  retention_period = 7
  backup_model = "logical"
  backup_time ="02:00–06:00"
}
```

##Argument Reference

The following arguments are supported:

- `mysql_id` - (Required) Instance ID to which policies will be applied. 

- `retention_period` - (Optional) Instance backup retention days. Valid values: [7-730]. And default value is 7.

- `backup_model` – (Optional) Backup method. Supported values include: physical - physical backup, and logical - logical backup.

- `backup_time` – (Optional) Instance backup time, in the format of "HH:mm-HH:mm". Time setting interval is four hours. Default to "02:00-06:00". The following value can be supported: 02:00\-06:00, 06:00\-10:00, 10:00\-14:00, 14:00\-18:00, 18:00\-22:00, and 22:00\-02:00.

##Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `binlog_period` – Retention period for binlog in days.