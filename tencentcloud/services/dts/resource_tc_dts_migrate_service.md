Provides a resource to create a dts migrate_service

Example Usage

```hcl
resource "tencentcloud_dts_migrate_service" "migrate_service" {
  src_database_type = "mysql"
  dst_database_type = "cynosdbmysql"
  src_region = "ap-guangzhou"
  dst_region = "ap-guangzhou"
  instance_class = "small"
  job_name = "tf_test_migration_job"
  tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
}

```
Import

dts migrate_service can be imported using the id, e.g.
```
$ terraform import tencentcloud_dts_migrate_service.migrate_service migrateService_id
```