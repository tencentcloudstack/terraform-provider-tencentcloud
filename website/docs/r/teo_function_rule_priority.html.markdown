---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_rule_priority"
sidebar_current: "docs-tencentcloud-resource-teo_function_rule_priority"
description: |-
  Provides a resource to create a teo teo_function_rule_priority
---

# tencentcloud_teo_function_rule_priority

Provides a resource to create a teo teo_function_rule_priority

## Example Usage

```hcl
resource "tencentcloud_teo_function_rule_priority" "teo_function_rule_priority" {
  function_id = "ef-txx7fnua"
  rule_ids = [
    "rule-equpbht3",
    "rule-ax28n3g6",
  ]
  zone_id = "zone-2qtuhspy7cr6"
}
```

## Argument Reference

The following arguments are supported:

* `function_id` - (Required, String, ForceNew) ID of the Function.
* `rule_ids` - (Required, List: [`String`]) he list of rule IDs. It is required to include all rule IDs after adjusting their priorities. The execution order of multiple rules follows a top-down sequence. If not specified, the original priority order will be maintained.
* `zone_id` - (Required, String, ForceNew) ID of the site.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo teo_function_rule_priority can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function_rule_priority.teo_function_rule_priority zone_id#function_id
```

