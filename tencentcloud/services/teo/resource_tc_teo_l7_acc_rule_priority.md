Provides a resource to create a teo l7_acc_rule_priority

Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rule_priority" "teo_l7_acc_rule_priority" {
    zone_id       = "zone-36bjhygh1bxe"
    rule_ids = [
        "rule-39pkyiu08edu",
        "rule-39pky6n21mkf",
    ]
}

```
Import

teo l7_acc_rule_priority can be imported using the zone_id, e.g.
````
terraform import tencentcloud_teo_l7_acc_rule_priority.teo_l7_acc_rule_priority zone-297z8rf93cfw
````