Use this data source to query security groups associated with TencentCloud PostgreSQL instance or read-only group.

Example Usage

Query security groups associated with a PostgreSQL instance

```hcl
data "tencentcloud_postgresql_db_instance_security_groups" "example" {
  db_instance_id = "postgres-ckwcgdf1"
}
```

Query security groups associated with a read-only group

```hcl
data "tencentcloud_postgresql_db_instance_security_groups" "example" {
  read_only_group_id = "pgrogrp-6fexzcmy"
}
```
