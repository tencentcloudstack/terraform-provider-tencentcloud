Provides a resource to create a tke kubernetes_native_node_pool

Example Usage

```hcl
resource "tencentcloud_kubernetes_native_node_pool" "kubernetes_native_node_pool" {
  labels = {
  }
  taints = {
  }
  tags = {
    tags = {
    }
  }
}
```

Import

tke kubernetes_native_node_pool can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_native_node_pool.kubernetes_native_node_pool kubernetes_native_node_pool_id
```
