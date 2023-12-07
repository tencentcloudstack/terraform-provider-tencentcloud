Provides a resource to create a mariadb security_groups

~> **NOTE:** If you use this resource, please do not use security_group_ids in tencentcloud_mariadb_instance resource

Example Usage

```hcl
resource "tencentcloud_mariadb_security_groups" "security_groups" {
  instance_id       = "tdsql-4pzs5b67"
  security_group_id = "sg-7kpsbxdb"
  product           = "mariadb"
}

```
Import

mariadb security_groups can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_security_groups.security_groups tdsql-4pzs5b67#sg-7kpsbxdb#mariadb
```