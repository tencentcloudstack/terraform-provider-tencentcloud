Use this data source to query detailed information of PostgreSQL db instance versions

Example Usage

Query all versions

```hcl
data "tencentcloud_postgresql_db_instance_versions" "example" {}
```

Query versions by storage type

```hcl
data "tencentcloud_postgresql_db_instance_versions" "example" {
  storage_type = "CLOUD_HSSD"
}
```
