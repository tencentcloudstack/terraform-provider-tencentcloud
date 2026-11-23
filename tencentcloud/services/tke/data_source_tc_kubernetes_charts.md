Use this data source to query detailed information of kubernetes cluster charts.

Example Usage

Query all kubernetes charts

```hcl
data "tencentcloud_kubernetes_charts" "example" {}
```

Query kubernetes charts by filter

```hcl
data "tencentcloud_kubernetes_charts" "example" {
  kind         = "network"
  arch         = "amd64"
  cluster_type = "tke"
}
```