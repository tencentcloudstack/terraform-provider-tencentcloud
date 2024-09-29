Provides a resource to create a teo teo_function_rule

Example Usage

```hcl
resource "tencentcloud_teo_function_rule" "teo_function_rule" {
    function_id   = "ef-txx7fnua"
    remark        = "aaa"
    zone_id       = "zone-2qtuhspy7cr6"

    function_rule_conditions {
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "host"
            values      = [
                "aaa.makn.cn",
            ]
        }
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "extension"
            values      = [
                ".txt",
            ]
        }
    }
    function_rule_conditions {
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "notequal"
            target      = "host"
            values      = [
                "aaa.makn.cn",
            ]
        }
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "extension"
            values      = [
                ".png",
            ]
        }
    }
}
```

Import

teo teo_function_rule can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function_rule.teo_function_rule zone_id#function_id#rule_id
```
