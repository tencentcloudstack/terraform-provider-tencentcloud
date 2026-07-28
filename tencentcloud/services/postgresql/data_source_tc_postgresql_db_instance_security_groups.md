Use this data source to query security groups associated with TencentCloud PostgreSQL instance or read-only group.

Example Usage

```hcl
data "tencentcloud_postgresql_db_instance_security_groups" "example" {
  db_instance_id = "postgres-gzg9jb2n"
}
```

```hcl
data "tencentcloud_postgresql_db_instance_security_groups" "example" {
  read_only_group_id = "pgro-xxxxx"
}
```
