Use this data source to query detailed information of postgresql default_parameters

Example Usage

```hcl
data "tencentcloud_postgresql_default_parameters" "default_parameters" {
  db_major_version = "13"
  db_engine = "postgresql"
}
```