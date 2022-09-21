---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_rule_engine"
sidebar_current: "docs-tencentcloud-resource-teo_rule_engine"
description: |-
  Provides a resource to create a teo rule_engine
---

# tencentcloud_teo_rule_engine

Provides a resource to create a teo rule_engine

## Example Usage

```hcl
resource "tencentcloud_teo_rule_engine" "rule_engine" {
  rule_name = "test-rule3"
  status    = "enable"
  tags      = {}
  zone_id   = "zone-297z8rf93cfw"

  rules {
    actions {

      normal_action {
        action = "Http2"

        parameters {
          name = "Switch"
          values = [
            "off",
          ]
        }
      }
    }
    actions {

      normal_action {
        action = "ForceRedirect"

        parameters {
          name = "Switch"
          values = [
            "on",
          ]
        }
        parameters {
          name = "RedirectStatusCode"
          values = [
            "302",
          ]
        }
      }
    }

    or {
      and {
        operator = "equal"
        target   = "host"
        values = [
          "www.toutiao2.com",
        ]
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `rule_name` - (Required, String) Rule name.
* `rules` - (Required, List) Rule items list.
* `status` - (Required, String) Status of the rule, valid value can be `enable` or `disable`.
* `zone_id` - (Required, String) Site ID.

The `rules` object supports the following:

* `or` - (Required, List) OR Conditions list of the rule. Rule would be triggered if any of the condition is true.
* `actions` - (Required, List) Actions list of the rule. See details in data source `rule_engine_setting`.

The `or` object supports the following:

* `and` - (Required, List) AND Conditions list of the rule. Rule would be triggered if all conditions are true.

The `and` object supports the following:

* `operator` - (Required, String) Condition operator. Valid values are `equal`, `notequal`.
* `target` - (Required, String) Condition target. Valid values:- `host`: Host of the URL.- `filename`: filename of the URL.- `extension`: file extension of the URL.- `full_url`: full url.- `url`: path of the URL.
* `values` - (Required, Set) Condition Value.

The `actions` object supports the following:

* `normal_action` - (Optional, List) Define a normal action.
* `code_action` - (Optional, List) Define a code action.
* `rewrite_action` - (Optional, List) Define a rewrite action.

The `normal_action` object supports the following:

* `action` - (Required, String) Action name.
* `parameters` - (Required, List) Action parameters.

The `parameters` object for `normal_action` supports the following:

* `name` - (Required, String) Parameter Name.
* `values` - (Required, Set) Parameter Values.

The `code_action` object supports the following:

* `action` - (Required, String) Action name.
* `parameters` - (Required, List) Action parameters.

The `parameters` object for `code_action` supports the following:

* `name` - (Required, String) Parameter Name.
* `status_code` - (Required, Int) HTTP status code to use.
* `values` - (Required, Set) Parameter Values.

The `rewrite_action` object supports the following:

* `action` - (Required, String) Action name.
* `parameters` - (Required, List) Action parameters.

The `parameters` object for `rewrite_action` supports the following:

* `action` - (Required, String) Action to take on the HEADER. Valid values: `add`, `del`, `set`.
* `name` - (Required, String) Target HEADER name.
* `values` - (Required, Set) Parameter Value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


## Import

teo rule_engine can be imported using the id#rule_id, e.g.
```
$ terraform import tencentcloud_teo_rule_engine.rule_engine zone-297z8rf93cfw#rule-ajol584a
```

