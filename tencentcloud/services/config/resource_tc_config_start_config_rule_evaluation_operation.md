Provides a resource to trigger a Config rule evaluation (one-shot operation).

Example Usage

Trigger evaluation by rule ID

```hcl
resource "tencentcloud_config_start_config_rule_evaluation_operation" "example" {
  rule_id = "cr-xhsd76j603v0a8ma0i73"
}
```

Trigger evaluation by compliance pack ID

```hcl
resource "tencentcloud_config_start_config_rule_evaluation_operation" "example" {
  compliance_pack_id = "cp-3kr5im1ssbg6tdo5jbi9"
}
```
