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
resource "tencentcloud_mysql_backup_policy" "default" {
  mysql_id         = "cdb-dnqksd9f"
  retention_period = 7
  backup_model     = "physical"
  backup_time      = "02:00-06:00"
}
```

## Argument Reference

The following arguments are supported:

* `mysql_id` - (Required, String, ForceNew) Instance ID to which policies will be applied.
* `backup_model` - (Optional, String) Backup method. Supported values include: `physical` - physical backup.
* `backup_time` - (Optional, String) Instance backup time, in the format of 'HH:mm-HH:mm'. Time setting interval is four hours. Default to `02:00-06:00`. The following value can be supported: `02:00-06:00`, `06:00-10:00`, `10:00-14:00`, `14:00-18:00`, `18:00-22:00`, and `22:00-02:00`.
* `retention_period` - (Optional, Int) Instance backup retention days. Valid value ranges: [7~730]. And default value is `7`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `binlog_period` - Retention period for binlog in days.


