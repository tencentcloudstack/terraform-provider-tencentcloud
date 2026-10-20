Provides a resource to create a FWM strategy

Example Usage

```hcl
resource "tencentcloud_fwm_strategy" "example" {
  product         = "cfm"
  receive_account = ["10000xxx"]

  pre_strategy {
    group_id = "rule-group-id-1"
    sequence = 1
  }

  post_strategy {
    group_id = "rule-group-id-2"
    sequence = 1
  }

  receive_group = ["group-xxx"]
}
```

Import

FWM strategy can be imported using the joined StrategyIds (joined by `#`), e.g. the create response returns StrategyIds `["s1","s2"]`, then import with `s1#s2`, e.g.

```
terraform import tencentcloud_fwm_strategy.example s1#s2
```