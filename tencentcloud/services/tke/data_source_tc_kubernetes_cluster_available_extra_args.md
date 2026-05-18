Use this data source to query the available custom extra arguments for TKE cluster components.

Example Usage

Query available extra args for a managed cluster

```hcl
data "tencentcloud_kubernetes_cluster_available_extra_args" "example" {
  cluster_version = "1.34.1"
  cluster_type    = "MANAGED_CLUSTER"
}
```

Query available extra args for an independent cluster

```hcl
data "tencentcloud_kubernetes_cluster_available_extra_args" "example" {
  cluster_version = "1.30.0"
  cluster_type    = "INDEPENDENT_CLUSTER"
}
```
