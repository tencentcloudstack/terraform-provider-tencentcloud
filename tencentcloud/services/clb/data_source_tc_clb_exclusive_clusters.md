Use this data source to query detailed information of clb exclusive_clusters

Example Usage

```hcl
data "tencentcloud_clb_exclusive_clusters" "exclusive_clusters" {
  filters {
    name = "zone"
    values = ["ap-guangzhou-1"]
  }
}
```