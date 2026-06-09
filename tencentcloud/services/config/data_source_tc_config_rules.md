Use this data source to query detailed information of Config rules.

Example Usage

Query all config rules

```hcl
data "tencentcloud_config_rules" "example" {}
```

Query config rules by name

```hcl
data "tencentcloud_config_rules" "example" {
  rule_name = "cam-user-mfa-check"
}
```

Query config rules by filters

```hcl
data "tencentcloud_config_rules" "example" {
  risk_level         = [1, 2]
  state              = "ACTIVE"
  compliance_result  = ["NON_COMPLIANT"]
  order_type         = "desc"
}
```
