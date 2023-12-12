Provides a resource to create a clickhouse backup strategy

Example Usage

```hcl
resource "tencentcloud_clickhouse_backup" "backup" {
	instance_id = "cdwch-xxxxxx"
	cos_bucket_name = "xxxxxx"
}

resource "tencentcloud_clickhouse_backup_strategy" "backup_strategy" {
  instance_id = "cdwch-xxxxxx"
  data_backup_strategy {
    week_days = "3"
    retain_days = 2
    execute_hour = 1
    back_up_tables {
      database = "iac"
      table = "my_table"
      total_bytes = 0
      v_cluster = "default_cluster"
      ips = "10.0.0.35"
    }
  }
  meta_backup_strategy {
	week_days = "1"
	retain_days = 2
	execute_hour = 3
  }
}
```

Import

clickhouse backup_strategy can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_backup_strategy.backup_strategy instance_id
```