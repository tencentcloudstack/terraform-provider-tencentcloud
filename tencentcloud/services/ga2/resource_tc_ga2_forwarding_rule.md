Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) layer-7 forwarding rule.

Example Usage

```hcl
resource "tencentcloud_ga2_forwarding_rule" "example" {
  global_accelerator_id = "ga-fhhs8w84"
  listener_id           = "lsr-dyy8jhzp"
  forwarding_policy_id  = "dm-rjssxr8k"

  rule_conditions {
    rule_condition_type  = "Path"
    rule_condition_value = ["/path"]
  }

  rule_actions {
    rule_action_type  = "ForwardGroup"
    rule_action_value = "epg-nt4iwozo"
  }
}
```

with origin headers

```hcl
resource "tencentcloud_ga2_forwarding_rule" "example1" {
  global_accelerator_id = "ga-fhhs8w84"
  listener_id           = "lsr-dyy8jhzp"
  forwarding_policy_id  = "dm-rjssxr8k"
  origin_host           = "2.2.2.2"

  rule_conditions {
    rule_condition_type  = "Path"
    rule_condition_value = ["/path"]
  }

  rule_actions {
    rule_action_type  = "ForwardGroup"
    rule_action_value = "epg-nt4iwozo"
  }

  origin_headers {
    key   = "key"
    value = "value"
  }
}
```

Import

GA2 forwarding rule can be imported using the composite id `<global_accelerator_id>#<listener_id>#<forwarding_policy_id>#<forwarding_rule_id>`, e.g.

```
terraform import tencentcloud_ga2_forwarding_rule.example ga-fhhs8w84#lsr-dyy8jhzp#dm-rjssxr8k#rule-757r3bk2
```
