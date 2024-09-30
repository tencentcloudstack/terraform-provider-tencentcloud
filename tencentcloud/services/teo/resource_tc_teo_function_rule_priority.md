Provides a resource to create a teo teo_function_rule_priority

Example Usage

```hcl
resource "tencentcloud_teo_function_rule_priority" "teo_function_rule_priority" {
    function_id = "ef-txx7fnua"
    rule_ids    = [
        "rule-equpbht3",
        "rule-ax28n3g6",
    ]
    zone_id     = "zone-2qtuhspy7cr6"
}
```

Import

teo teo_function_rule_priority can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function_rule_priority.teo_function_rule_priority zone_id#function_id
```
