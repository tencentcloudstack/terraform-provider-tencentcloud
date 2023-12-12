Use this data source to query detailed information of clb cluster_resources

Example Usage

```hcl
data "tencentcloud_clb_cluster_resources" "cluster_resources" {
  filters {
    name = "idle"
    values = ["True"]
  }
}
```