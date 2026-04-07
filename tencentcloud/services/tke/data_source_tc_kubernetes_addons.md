Use this data source to query detailed information of kubernetes addons.

Example Usage

```hcl
data "tencentcloud_kubernetes_addons" "example" {
  cluster_id = "cls-5yezvaxo"
}
```

Or

```hcl
data "tencentcloud_kubernetes_addons" "example" {
  cluster_id = "cls-5yezvaxo"
  addon_name = "ip-masq-agent"
}
```
