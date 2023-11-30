Use this data source to query detailed information of kubernetes cluster_instances

Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_instances" "cluster_instances" {
  cluster_id    = "cls-ely08ic4"
  instance_ids  = ["ins-kqmx8dm2"]
  instance_role = "WORKER"
  filters {
    name   = "nodepool-id"
    values = ["np-p4e6whqu"]
  }
}
```