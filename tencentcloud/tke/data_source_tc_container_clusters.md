Get container clusters in the current region.

Use this data source to get container clusters in the current region. By default every clusters in current region will be returned.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_kubernetes_clusters.

Example Usage

```hcl
data "tencentcloud_container_clusters" "foo" {
}
```