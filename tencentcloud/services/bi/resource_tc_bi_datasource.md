Provides a resource to create a bi datasource

Example Usage

```hcl
resource "tencentcloud_bi_datasource" "datasource" {
  charset     = "utf8"
  db_host     = "bj-cdb-1lxqg5r6.sql.tencentcdb.com"
  db_name     = "tf-test"
  db_port     = 63694
  db_type     = "MYSQL"
  db_pwd      = "ABc123,,,"
  db_user     = "root"
  project_id  = 11015030
  source_name = "tf-source-name"
}
```

Import

bi datasource can be imported using the id, e.g.

```
terraform import tencentcloud_bi_datasource.datasource datasource_id
```