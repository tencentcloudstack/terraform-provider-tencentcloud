Use this data source to query detailed information of tsf group_gateways

Example Usage

```hcl
data "tencentcloud_tsf_group_gateways" "group_gateways" {
  gateway_deploy_group_id = "group-aeoej4qy"
  search_word = "test"
}
```