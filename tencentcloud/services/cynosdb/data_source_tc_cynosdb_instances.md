Use this data source to query detailed information of Cynosdb instances.

Example Usage

```hcl
data "tencentcloud_cynosdb_instances" "foo" {
  instance_id   = "cynosdbmysql-ins-0wln9u6w"
  project_id    = 0
  db_type       = "MYSQL"
  instance_name = "test"
}
```