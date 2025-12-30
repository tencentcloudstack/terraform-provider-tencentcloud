Use this data source to query PostgreSQL instances

Example Usage

Query all postgresql instances

```hcl
data "tencentcloud_postgresql_instances" "example" {}
```

Query postgresql instances by filters

```hcl
data "tencentcloud_postgresql_instances" "example" {
  id         = "postgres-gngyhl9d"
  name       = "tf-example"
  project_id = "1235143"
}
```
