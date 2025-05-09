Provides a resource to create a cynosdb ssl

Example Usage

```hcl
resource "tencentcloud_cynosdb_ssl" "cynosdb_ssl" {
  cluster_id = "cynosdbmysql-1e0nzayx"
  instance_id = "cynosdbmysql-ins-pfsv6q1e"
  status = "ON"
}
```