Provides a resource to create a TEO function rule

Example Usage

```hcl
resource "tencentcloud_teo_function_rule" "example" {
    function_id   = "ef-m01xn26e"
    remark        = "remark."
    trigger_type  = "direct"
    zone_id       = "zone-3fkff38fyw8s"

    function_rule_conditions {
        rule_conditions {
            ignore_case = false
            operator    = "equal"
            target      = "host"
            values      = [
                "test.makn.cn",
            ]
        }
        rule_conditions {
            ignore_case = false
            operator    = "equal"
            target      = "url"
            values      = [
                "/path",
            ]
        }
    }
}
```

Import

teo teo_function_rule can be imported using the zoneId#functionId#ruleId, e.g.

```
terraform import tencentcloud_teo_function_rule.example zone-3fkff38fyw8s#ef-m01xn26e#rule-yuvufj6h
```
