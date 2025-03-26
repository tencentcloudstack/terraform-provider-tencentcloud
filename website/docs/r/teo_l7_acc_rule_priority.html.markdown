---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_l7_acc_rule_priority"
sidebar_current: "docs-tencentcloud-resource-teo_l7_acc_rule_priority"
description: |-
  Provides a resource to create a teo l7_acc_rule_priority
---

# tencentcloud_teo_l7_acc_rule_priority

Provides a resource to create a teo l7_acc_rule_priority

## Example Usage

```hcl
resource "tencentcloud_teo_l7_acc_rule_priority" "teo_l7_acc_rule_priority" {
  zone_id = "zone-36bjhygh1bxe"
  rule_ids = [
    "rule-39pkyiu08edu",
    "rule-39pky6n21mkf",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `rule_ids` - (Required, List: [`String`]) The final priority order of the rule ID list will be adjusted to the order of the rule ID list, and will be executed from the front to the back. The later rules will overwrite the earlier rules.
* `zone_id` - (Required, String, ForceNew) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo l7_acc_rule_priority can be imported using the zone_id, e.g.
````
terraform import tencentcloud_teo_l7_acc_rule_priority.teo_l7_acc_rule_priority zone-297z8rf93cfw
````

