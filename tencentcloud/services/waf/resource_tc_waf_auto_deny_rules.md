Provides a resource to create a waf auto_deny_rules

Example Usage

```hcl
resource "tencentcloud_waf_auto_deny_rules" "example" {
  domain              = "demo.waf.com"
  attack_threshold    = 20
  time_threshold      = 12
  deny_time_threshold = 5
}
```

Import

waf auto_deny_rules can be imported using the id, e.g.

```
terraform import tencentcloud_waf_auto_deny_rules.example demo.waf.com
```