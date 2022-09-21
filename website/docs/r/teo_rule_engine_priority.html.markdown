---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_rule_engine_priority"
sidebar_current: "docs-tencentcloud-resource-teo_rule_engine_priority"
description: |-
  Provides a resource to create a teo rule_engine_priority
---

# tencentcloud_teo_rule_engine_priority

Provides a resource to create a teo rule_engine_priority

## Example Usage

```hcl
resource "tencentcloud_teo_rule_engine_priority" "rule_engine_priority" {
  zone_id = "zone-294v965lwmn6"

  rules_priority {
    index = 0
    value = "rule-m9jlttua"
  }
  rules_priority {
    index = 1
    value = "rule-m5l9t4k1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Site ID.
* `rules_priority` - (Optional, List) Priority of rules.

The `rules_priority` object supports the following:

* `index` - (Optional, Int) Priority order of rules.
* `value` - (Optional, String) Priority of rules id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo rule_engine_priority can be imported using the zone_id, e.g.
```
$ terraform import tencentcloud_teo_rule_engine_priority.rule_engine_priority zone-294v965lwmn6
```

