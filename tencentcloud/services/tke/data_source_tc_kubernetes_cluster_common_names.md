Provide a datasource to query cluster CommonNames.

Example Usage

Query common names by subaccount uins

```hcl
data "tencentcloud_kubernetes_cluster_common_names" "example" {
  cluster_id      = "cls-fdy7hm1q"
  subaccount_uins = ["100037718139", "100031340176"]
}
```

Query common names by role ids

```hcl
data "tencentcloud_kubernetes_cluster_common_names" "example" {
  cluster_id = "cls-fdy7hm1q"
  role_ids   = ["4611686018441060141"]
}
```