Use this data source to query detailed information of IGTM strategy list

Example Usage

```hcl
data "tencentcloud_igtm_strategy_list" "example" {
  instance_id = "gtm-uukztqtoaru"
  filters {
    name  = "StrategyName"
    value = ["tf-example"]
    fuzzy = true
  }
}
```
