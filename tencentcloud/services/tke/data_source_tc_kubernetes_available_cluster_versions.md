Use this data source to query detailed information of kubernetes available_cluster_versions

Example Usage

```hcl
data "tencentcloud_kubernetes_available_cluster_versions" "query_by_id" {
  cluster_id = "xxx"
}

output "versions_id"{
  description = "Query versions from id."
  value = data.tencentcloud_kubernetes_available_cluster_versions.query_by_id.versions
}

data "tencentcloud_kubernetes_available_cluster_versions" "query_by_ids" {
  cluster_ids = ["xxx"]
}

output "versions_ids"{
  description = "Query versions from ids."
  value = data.tencentcloud_kubernetes_available_cluster_versions.query_by_ids.clusters
}
```