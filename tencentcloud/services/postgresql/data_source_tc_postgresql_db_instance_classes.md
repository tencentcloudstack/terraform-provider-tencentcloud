Use this data source to query detailed information of PostgreSQL db instance classes

Example Usage

```hcl
data "tencentcloud_postgresql_db_instance_classes" "example" {
  zone             = "ap-guangzhou-7"
  db_engine        = "postgresql"
  db_major_version = "13"
  storage_type     = "CLOUD_HSSD"
}
```
