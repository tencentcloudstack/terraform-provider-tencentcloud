Use this data source to query detailed information of Config rule evaluation results (rule dimension).

Example Usage

Query evaluation results by rule ID

```hcl
data "tencentcloud_config_rule_evaluation_results" "example" {
  config_rule_id = "cr-pHmVQS1UpihV4MSTkmIo"
}
```

Query evaluation results with filters

```hcl
data "tencentcloud_config_rule_evaluation_results" "example" {
  config_rule_id  = "cr-pHmVQS1UpihV4MSTkmIo"
  compliance_type = ["NON_COMPLIANT"]
  resource_type   = ["QCS::CVM::Instance"]
}
```
