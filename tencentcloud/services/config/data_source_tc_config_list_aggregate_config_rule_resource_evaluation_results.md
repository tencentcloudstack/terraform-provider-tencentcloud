Use this data source to query Config aggregate config rule resource evaluation results.

Example Usage

Query aggregate config rule resource evaluation results for a specific account group

```hcl
data "tencentcloud_config_list_aggregate_config_rule_resource_evaluation_results" "example" {
  account_group_id = "ca-aggregator-abcdefg"
}
```

Query aggregate config rule resource evaluation results for all account groups

```hcl
data "tencentcloud_config_list_aggregate_config_rule_resource_evaluation_results" "example" {}
```