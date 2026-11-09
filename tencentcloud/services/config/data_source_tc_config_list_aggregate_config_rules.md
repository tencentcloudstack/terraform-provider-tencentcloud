Use this data source to query the list of Config aggregate config rules under an account group.

Example Usage

Query aggregate config rules by account group id

```hcl
data "tencentcloud_config_list_aggregate_config_rules" "example" {
  account_group_id = "ca-aggregator-abcdefg"
}
```

Query aggregate config rules with filters

```hcl
data "tencentcloud_config_list_aggregate_config_rules" "example" {
  account_group_id  = "ca-aggregator-abcdefg"
  risk_level        = [1, 2]
  state             = "ACTIVE"
  compliance_result = ["NON_COMPLIANT"]
  order_type        = "desc"
}
```