Use this data source to query detailed information of kubernetes cluster charts.

Example Usage

Query all kubernetes charts

```hcl
data "tencentcloud_kubernetes_charts" "example" {}
```

Query kubernetes charts by filter

```hcl
data "tencentcloud_kubernetes_charts" "example" {
  kind = "network"
}

data "tencentcloud_kubernetes_charts" "example" {
  arch = "amd64"
}

data "tencentcloud_kubernetes_charts" "example" {
  cluster_type = "tke"
}
```