Provides a resource to create a cynosdb upgrade_proxy_version

Example Usage

```hcl
resource "tencentcloud_cynosdb_upgrade_proxy_version" "upgrade_proxy_version" {
  cluster_id = "cynosdbmysql-bws8h88b"
  dst_proxy_version = "1.3.7"
}
```