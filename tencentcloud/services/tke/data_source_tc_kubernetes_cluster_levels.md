Provide a datasource to query TKE cluster levels.

Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_levels" "foo" {}

output "level5" {
	value = data.tencentcloud_kubernetes_cluster_levels.foo.list.0.alias
}
```