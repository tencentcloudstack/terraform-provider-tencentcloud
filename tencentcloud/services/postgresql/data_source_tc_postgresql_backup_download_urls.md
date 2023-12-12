Use this data source to query detailed information of postgresql backup_download_urls

Example Usage

```hcl
data "tencentcloud_postgresql_log_backups" "log_backups" {
	min_finish_time = "%s"
	max_finish_time = "%s"
	filters {
		  name = "db-instance-id"
		  values = [local.pgsql_id]
	}
	order_by = "StartTime"
	order_by_type = "desc"

  }

data "tencentcloud_postgresql_backup_download_urls" "backup_download_urls" {
  db_instance_id = local.pgsql_id
  backup_type = "LogBackup"
  backup_id = data.tencentcloud_postgresql_log_backups.log_backups.log_backup_set.0.id
  url_expire_time = 12
  backup_download_restriction {
		restriction_type = "NONE"
		vpc_restriction_effect = "ALLOW"
		vpc_id_set = [local.vpc_id]
		ip_restriction_effect = "ALLOW"
		ip_set = ["0.0.0.0"]
  }
}
```