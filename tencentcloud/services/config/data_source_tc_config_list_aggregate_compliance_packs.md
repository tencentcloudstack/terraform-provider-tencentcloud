Use this data source to query the list of Config aggregate compliance packs by account group.

Example Usage

Query all aggregate compliance packs of an account group

```hcl
data "tencentcloud_config_list_aggregate_compliance_packs" "example" {
  account_group_id = "ca-Zp9EiU7Pa4Vy0lFFm8XM"
}
```

Query aggregate compliance packs with filters

```hcl
data "tencentcloud_config_list_aggregate_compliance_packs" "example" {
  account_group_id     = "ca-Zp9EiU7Pa4Vy0lFFm8XM"
  compliance_pack_name = "tf-example"
  risk_level           = [1, 2]
  status               = "UN_ACTIVE"
  compliance_result    = ["NON_COMPLIANT"]
  order_type           = "desc"
}
```
