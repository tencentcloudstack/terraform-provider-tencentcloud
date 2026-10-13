---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_log_file_retention_period"
sidebar_current: "docs-tencentcloud-resource-mariadb_log_file_retention_period"
description: |-
  Provides a resource to create a mariadb log_file_retention_period
---

# tencentcloud_mariadb_log_file_retention_period

Provides a resource to create a mariadb log_file_retention_period

## Example Usage

```hcl
resource "tencentcloud_mariadb_log_file_retention_period" "log_file_retention_period" {
  instance_id = "tdsql-4pzs5b67"
  days        = "8"
}
```

## Argument Reference

The following arguments are supported:

* `days` - (Required, Int) The number of days to save, cannot exceed 30.
* `instance_id` - (Required, String) instance id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyLogFileRetentionPeriod](https://cloud.tencent.com/document/api/237/16151) - 修改备份日志保存天数
* Read: [DescribeLogFileRetentionPeriod](https://cloud.tencent.com/document/api/237/16152) - 查看备份日志备份天数
* Update: [ModifyLogFileRetentionPeriod](https://cloud.tencent.com/document/api/237/16151) - 修改备份日志保存天数

## Import

mariadb log_file_retention_period can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_log_file_retention_period.log_file_retention_period tdsql-4pzs5b67
```

