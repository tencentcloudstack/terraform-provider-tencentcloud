Provides a resource to create a cynosdb ssl

Example Usage

```hcl
resource "tencentcloud_cynosdb_ssl" "example" {
  cluster_id  = "cynosdbmysql-1e0nzayx"
  instance_id = "cynosdbmysql-ins-pfsv6q1e"
  status      = "ON"
}
```

Import

cynosdb ssl can be imported using the clusterId#instanceId, e.g.

```
terraform import tencentcloud_cynosdb_ssl.example cynosdbmysql-1e0nzayx#cynosdbmysql-ins-pfsv6q1e
```