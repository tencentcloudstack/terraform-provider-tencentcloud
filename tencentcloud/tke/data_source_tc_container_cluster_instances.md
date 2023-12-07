Get all instances of the specific cluster.

Use this data source to get all instances in a specific cluster.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_kubernetes_clusters.

Example Usage

```hcl
data "tencentcloud_container_cluster_instances" "foo_instance" {
  cluster_id = "cls-abcdefg"
}
```