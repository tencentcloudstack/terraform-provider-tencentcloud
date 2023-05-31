---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_backup_download_urls"
sidebar_current: "docs-tencentcloud-datasource-postgresql_backup_download_urls"
description: |-
  Use this data source to query detailed information of postgresql backup_download_urls
---

# tencentcloud_postgresql_backup_download_urls

Use this data source to query detailed information of postgresql backup_download_urls

## Example Usage

```hcl
data "tencentcloud_postgresql_log_backups" "log_backups" {
  min_finish_time = "%s"
  max_finish_time = "%s"
  filters {
    name   = "db-instance-id"
    values = [local.pgsql_id]
  }
  order_by      = "StartTime"
  order_by_type = "desc"

}

data "tencentcloud_postgresql_backup_download_urls" "backup_download_urls" {
  db_instance_id  = local.pgsql_id
  backup_type     = "LogBackup"
  backup_id       = data.tencentcloud_postgresql_log_backups.log_backups.log_backup_set.0.id
  url_expire_time = 12
  backup_download_restriction {
    restriction_type       = "NONE"
    vpc_restriction_effect = "ALLOW"
    vpc_id_set             = [local.vpc_id]
    ip_restriction_effect  = "ALLOW"
    ip_set                 = ["0.0.0.0"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `backup_id` - (Required, String) Unique backup ID.
* `backup_type` - (Required, String) Backup type. Valid values: `LogBackup`, `BaseBackup`.
* `db_instance_id` - (Required, String) Instance ID.
* `backup_download_restriction` - (Optional, List) Backup download restriction.
* `result_output_file` - (Optional, String) Used to save results.
* `url_expire_time` - (Optional, Int) Validity period of a URL, which is 12 hours by default.

The `backup_download_restriction` object supports the following:

* `ip_restriction_effect` - (Optional, String) Whether IP is allowed. Valid values: `ALLOW` (allow), `DENY` (deny).
* `ip_set` - (Optional, Set) Whether it is allowed to download IP list of the backup files.
* `restriction_type` - (Optional, String) Type of the network restrictions for downloading backup files. Valid values: `NONE` (backups can be downloaded over both private and public networks), `INTRANET` (backups can only be downloaded over the private network), `CUSTOMIZE` (backups can be downloaded over specified VPCs or at specified IPs).
* `vpc_id_set` - (Optional, Set) Whether it is allowed to download the VPC ID list of the backup files.
* `vpc_restriction_effect` - (Optional, String) Whether VPC is allowed. Valid values: `ALLOW` (allow), `DENY` (deny).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backup_download_url` - Backup download URL.


