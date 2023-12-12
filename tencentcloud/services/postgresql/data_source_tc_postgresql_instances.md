Use this data source to query postgresql instances

Example Usage

```hcl
data "tencentcloud_postgresql_instances" "name" {
  name = "test"
}

data "tencentcloud_postgresql_instances" "project" {
  project_id = 0
}

data "tencentcloud_postgresql_instances" "id" {
  id = "postgres-h9t4fde1"
}
```