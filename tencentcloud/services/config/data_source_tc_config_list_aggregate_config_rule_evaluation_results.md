Use this data source to query Config aggregate config rule evaluation results (rule dimension) for an account group.

Example Usage

Query evaluation results for a config rule under an account group

```hcl
data "tencentcloud_config_list_aggregate_config_rule_evaluation_results" "example" {
  config_rule_id   = "cr-abcdefg1234"
  account_group_id = "ca-aggregator-abcdefg"
}
```

Query evaluation results with filters

```hcl
data "tencentcloud_config_list_aggregate_config_rule_evaluation_results" "example" {
  config_rule_id   = "cr-abcdefg1234"
  account_group_id = "ca-aggregator-abcdefg"
  compliance_type   = ["NON_COMPLIANT"]
  resource_type     = ["QCS::CVM::Instance"]
}
```