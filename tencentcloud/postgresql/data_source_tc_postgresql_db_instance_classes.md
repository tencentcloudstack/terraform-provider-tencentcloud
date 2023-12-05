Use this data source to query detailed information of postgresql db_instance_classes

Example Usage

```hcl
data "tencentcloud_postgresql_db_instance_classes" "db_instance_classes" {
  zone = "ap-guangzhou-7"
  db_engine = "postgresql"
  db_major_version = "13"
}
```