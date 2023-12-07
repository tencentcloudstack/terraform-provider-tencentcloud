Provides a resource to create a cynosdb reload_proxy_node

Example Usage

```hcl
resource "tencentcloud_cynosdb_reload_proxy_node" "reload_proxy_node" {
  cluster_id     = "cynosdbmysql-cgd2gpwr"
  proxy_group_id = "cynosdbmysql-proxy-8lqtl8pk"
}
```

Import

cynosdb reload_proxy_node can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_reload_proxy_node.reload_proxy_node reload_proxy_node_id
```