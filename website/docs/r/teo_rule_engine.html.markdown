---
subcategory: "Teo"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_rule_engine"
sidebar_current: "docs-tencentcloud-resource-teo_rule_engine"
description: |-
  Provides a resource to create a teo ruleEngine
---

# tencentcloud_teo_rule_engine

Provides a resource to create a teo ruleEngine

## Example Usage

```hcl
resource "tencentcloud_teo_rule_engine" "rule_engine" {
  zone_id   = tencentcloud_teo_zone.zone.id
  rule_name = "rule0"
  status    = "enable"

  rules {
    conditions {
      conditions {
        operator = "equal"
        target   = "host"
        values = [
          "www.sfurnace.work",
        ]
      }
    }

    actions {
      normal_action {
        action = "MaxAge"

        parameters {
          name = "FollowOrigin"
          values = [
            "on",
          ]
        }
        parameters {
          name = "MaxAgeTime"
          values = [
            "0",
          ]
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `rules` - (Required, List) Rule items list.
* `status` - (Required, String) Status of the rule, valid value can be `enable` or `disable`.
* `zone_id` - (Required, String) Site ID.
* `rule_name` - (Optional, String) Rule name.

The `actions` object supports the following:

* `code_action` - (Optional, List) Define a code action.
* `normal_action` - (Optional, List) Define a normal action.
* `rewrite_action` - (Optional, List) Define a rewrite action.

The `code_action` object supports the following:

* `action` - (Required, String) Action name.
* `parameters` - (Required, List) Action parameters.

The `conditions` object supports the following:

* `conditions` - (Required, List) AND Conditions list of the rule. Rule would be triggered if all conditions are true.

The `conditions` object supports the following:

* `operator` - (Required, String) Condition operator. Valid values are `equal`, `notequal`.
* `target` - (Required, String) Condition target. Valid values:- host: Host of the URL.- filename: filename of the URL.- extension: file extension of the URL.- full_url: full url.- url: path of the URL.
* `values` - (Required, Set) Condition Value.

The `normal_action` object supports the following:

* `action` - (Required, String) Action name.
* `parameters` - (Required, List) Action parameters.

The `parameters` object supports the following:

* `action` - (Required, String) Action to take on the HEADER.
* `name` - (Required, String) Target HEADER name.
* `values` - (Required, Set) Parameter Value.

The `parameters` object supports the following:

* `name` - (Required, String) Parameter Name.
* `status_code` - (Required, Int) HTTP status code to use.
* `values` - (Required, Set) Parameter Values.

The `parameters` object supports the following:

* `name` - (Required, String) Parameter Name.
* `values` - (Required, Set) Parameter Values.

The `rewrite_action` object supports the following:

* `action` - (Required, String) Action name.
* `parameters` - (Required, List) Action parameters.

The `rules` object supports the following:

* `actions` - (Required, List) Actions list of the rule. See details in data source `rule_engine_setting`.
* `conditions` - (Required, List) OR Conditions list of the rule. Rule would be triggered if any of the condition is true.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


## Import

teo ruleEngine can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_rule_engine.rule_engine zoneId#ruleId
```

