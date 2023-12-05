Provides a resource to create a cynosdb cluster_databases

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_databases" "cluster_databases" {
  cluster_id    = "cynosdbmysql-bws8h88b"
  db_name       = "terraform-test"
  character_set = "utf8"
  collate_rule  = "utf8_general_ci"
  user_host_privileges {
    db_user_name = "root"
    db_host      = "%"
    db_privilege = "READ_ONLY"
  }
  description = "terraform test"
}
```

Import

cynosdb cluster_databases can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_databases.cluster_databases cluster_databases_id
```