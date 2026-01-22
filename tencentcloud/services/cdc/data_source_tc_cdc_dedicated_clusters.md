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

data "tencentcloud_cdc_dedicated_clusters" "example" {
  dedicated_cluster_ids = [
    "cluster-aiaui7ei",
    "cluster-262n63e8"
  ]
}

data "tencentcloud_cdc_dedicated_clusters" "example" {
  zones = [
    "ap-guangzhou-2",
    "ap-guangzhou-6"
  ]
}

data "tencentcloud_cdc_dedicated_clusters" "example" {
  site_ids = [
    "site-2qu42ele",
    "site-fp8gp962"
  ]
}

data "tencentcloud_cdc_dedicated_clusters" "example" {
  lifecycle_statuses = [
    "PENDING",
    "RUNNING"
  ]
}
```
