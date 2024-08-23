Provides a resource to create a tke kubernetes_native_node_pools

Example Usage

```hcl
resource "tencentcloud_kubernetes_native_node_pools" "kubernetes_native_node_pools" {
  labels = {
  }
  taints = {
  }
  tags = {
    tags = {
    }
  }
  native = {
    system_disk = {
    }
    data_disks = {
    }
  }
}
```

Import

tke kubernetes_native_node_pools can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_native_node_pools.kubernetes_native_node_pools kubernetes_native_node_pools_id
```
