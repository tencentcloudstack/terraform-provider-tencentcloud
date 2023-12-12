Use this data source to query detailed information of kubernetes clusters.

Example Usage

```hcl
data "tencentcloud_kubernetes_clusters" "name" {
  cluster_name = "terraform"
}

data "tencentcloud_kubernetes_clusters" "id" {
  cluster_id = "cls-godovr32"
}
```