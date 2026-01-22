Use this data source to query detailed information of IGTM address pool list

Example Usage

Query all address pool list

```hcl
data "tencentcloud_igtm_address_pool_list" "example" {}
```

Query address pool list by filter

```hcl
data "tencentcloud_igtm_address_pool_list" "example" {
  filters {
    name  = "PoolName"
    value = ["tf-example"]
    fuzzy = true
  }
}
```
