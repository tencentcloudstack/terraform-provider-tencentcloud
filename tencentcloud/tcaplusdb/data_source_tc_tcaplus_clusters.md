Use this data source to query TcaplusDB clusters.

Example Usage

```hcl
data "tencentcloud_tcaplus_clusters" "name" {
  cluster_name = "cluster"
}
data "tencentcloud_tcaplus_clusters" "id" {
  cluster_id = tencentcloud_tcaplus_cluster.test.id
}
data "tencentcloud_tcaplus_clusters" "idname" {
  cluster_id   = tencentcloud_tcaplus_cluster.test.id
  cluster_name = "cluster"
}
```