Use this data source to query detailed information of CDC dedicated clusters

Example Usage

Query all dedicated clusters

```hcl
data "tencentcloud_cdc_dedicated_clusters" "example" {}
```

Query dedicated clusters by filters

```hcl
data "tencentcloud_cdc_dedicated_clusters" "example" {
  name = "tf-example"
}

data "tencentcloud_cdc_dedicated_clusters" "example_2" {
  dedicated_cluster_ids = [
    "cluster-aiaui7ei",
    "cluster-262n63e8"
  ]
}

data "tencentcloud_cdc_dedicated_clusters" "example_3" {
  zones = [
    "ap-guangzhou-2",
    "ap-guangzhou-6"
  ]
}

data "tencentcloud_cdc_dedicated_clusters" "example_4" {
  site_ids = [
    "site-2qu42ele",
    "site-fp8gp962"
  ]
}

data "tencentcloud_cdc_dedicated_clusters" "example_5" {
  lifecycle_statuses = [
    "PENDING",
    "RUNNING"
  ]
}
```
