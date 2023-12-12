Provides a resource to create a cynosdb wan

Example Usage

```hcl
resource "tencentcloud_cynosdb_wan" "wan" {
  cluster_id      = "cynosdbmysql-bws8h88b"
  instance_grp_id = "cynosdbmysql-grp-lxav0p9z"
}
```

Import

cynosdb wan can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_wan.wan cynosdbmysql-bws8h88b#cynosdbmysql-grp-lxav0p9z
```