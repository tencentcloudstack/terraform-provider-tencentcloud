Provides a resource to create a GA2 (Global Accelerator 2) forwarding rule.

Example Usage

```hcl
resource "tencentcloud_ga2_forwarding_rule" "example" {
  global_accelerator_id = "ga2-xxxxxxxx"
  listener_id           = "lis-xxxxxxxx"
  forwarding_policy_id  = "fp-xxxxxxxx"

  rule_conditions {
    rule_condition_type  = "HOST"
    rule_condition_value = ["example.com"]
  }

  rule_actions {
    rule_action_type  = "ForwardGroup"
    rule_action_value = "eg-xxxxxxxx"
  }

  origin_headers {
    key   = "X-Custom-Header"
    value = "custom-value"
  }

  enable_origin_sni = true
  origin_sni        = "example.com"
  origin_host       = "origin.example.com"
}
```

Import

GA2 forwarding rule can be imported using the composite ID `<global_accelerator_id>#<listener_id>#<forwarding_policy_id>#<forwarding_rule_id>`, e.g.

```
terraform import tencentcloud_ga2_forwarding_rule.example ga2-xxxxxxxx#lis-xxxxxxxx#fp-xxxxxxxx#fr-xxxxxxxx
```
