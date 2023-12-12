Provides a resource to create a redis backup_operation

Example Usage

Manually back up the Redis instance, and the backup data is kept for 7 days

```hcl
data "tencentcloud_mysql_instance" "foo" {}

resource "tencentcloud_redis_backup_operation" "backup_operation" {
  instance_id = data.tencentcloud_mysql_instance.foo.instance_list[0].mysql_id
  remark = "manually back"
  storage_days = 7
}
```