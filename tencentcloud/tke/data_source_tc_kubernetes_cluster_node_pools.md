Use this data source to query detailed information of kubernetes cluster_node_pools

Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_node_pools" "cluster_node_pools" {
  cluster_id = "cls-kzilgv5m"
  filters {
    name   = "NodePoolsName"
    values = ["mynodepool_xxxx"]
  }
  filters {
    name   = "NodePoolsId"
    values = ["np-ngjwhdv4"]
  }
}
```
