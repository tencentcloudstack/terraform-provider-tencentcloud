Provides a resource to create a ckafka acl_rule

Example Usage

```hcl
resource "tencentcloud_ckafka_acl_rule" "acl_rule" {
  instance_id = "ckafka-xxx"
  resource_type = "Topic"
  pattern_type = "PREFIXED"
  rule_name = "RuleName"
  rule_list {
		operation = "All"
		permission_type = "Deny"
		host = "*"
		principal = "User:*"

  }
  pattern = "prefix"
  is_applied = 1
}
```

Import

ckafka acl_rule can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_acl_rule.acl_rule acl_rule_id
```