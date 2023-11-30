Provides a resource to create a cynosdb security_group

Example Usage

```hcl
resource "tencentcloud_cynosdb_security_group" "test" {
  cluster_id = "cynosdbmysql-bws8h88b"
  security_group_ids = ["sg-baxfiao5"]
  instance_group_type = "RO"
}
```

Import

cynosdb security_group can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_security_group.security_group ${cluster_id}#${instance_group_type}
```