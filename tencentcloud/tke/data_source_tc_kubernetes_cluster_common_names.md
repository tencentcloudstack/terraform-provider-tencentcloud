Provide a datasource to query cluster CommonNames.

Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_common_names" "foo" {
  cluster_id = "cls-12345678"
  subaccount_uins = ["1234567890", "0987654321"]
}
```