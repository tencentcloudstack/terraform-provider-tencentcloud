---
subcategory: "EventBridge(EB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eb_event_rule"
sidebar_current: "docs-tencentcloud-resource-eb_event_rule"
description: |-
  Provides a resource to create a eb event_rule
---

# tencentcloud_eb_event_rule

Provides a resource to create a eb event_rule

## Example Usage

```hcl
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_event_rule" "event_rule" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_name    = "tf-event_rule"
  description  = "event rule desc"
  enable       = true
  event_pattern = jsonencode(
    {
      source = "apigw.cloud.tencent"
      type = [
        "connector:apigw",
      ]
    }
  )
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `event_bus_id` - (Required, String) event bus Id.
* `event_pattern` - (Required, String) Reference: [Event Mode](https://cloud.tencent.com/document/product/1359/56084).
* `rule_name` - (Required, String) Event rule name, which can only contain letters, numbers, underscores, hyphens, starts with a letter and ends with a number or letter, 2~60 characters.
* `description` - (Optional, String) Event set description, unlimited character type, description within 200 characters.
* `enable` - (Optional, Bool) Enable switch.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - event rule id.


## Import

eb event_rule can be imported using the id, e.g.

```
terraform import tencentcloud_eb_event_rule.event_rule event_rule_id
```

