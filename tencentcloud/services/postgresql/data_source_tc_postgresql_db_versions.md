Use this data source to query detailed information of postgres db_versions

Example Usage

Query all DB versions

```hcl
data "tencentcloud_postgresql_db_versions" "example" {}
```

Query DB versions by filters

```hcl
data "tencentcloud_postgresql_db_versions" "example" {
  db_version        = "16.0"
  db_major_version  = "16"
  db_kernel_version = "v16.0_r1.0"
}
```
