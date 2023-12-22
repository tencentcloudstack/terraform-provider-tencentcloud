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
resource "tencentcloud_teo_rule_engine" "rule1" {
  zone_id   = tencentcloud_teo_zone.example.id
  rule_name = "test-rule"
  status    = "disable"

  rules {
    actions {
      normal_action {
        action = "UpstreamUrlRedirect"
        parameters {
          name = "Type"
          values = [
            "Path",
          ]
        }
        parameters {
          name = "Action"
          values = [
            "addPrefix",
          ]
        }
        parameters {
          name = "Value"
          values = [
            "/sss",
          ]
        }
      }
    }

    or {
      and {
        operator    = "equal"
        target      = "host"
        ignore_case = false
        values = [
          "a.tf-teo-t.xyz",
        ]
      }
      and {
        operator    = "equal"
        target      = "extension"
        ignore_case = false
        values = [
          "jpg",
        ]
      }
    }
    or {
      and {
        operator    = "equal"
        target      = "filename"
        ignore_case = false
        values = [
          "test.txt",
        ]
      }
    }

    sub_rules {
      tags = ["png"]
      rules {
        or {
          and {
            operator    = "notequal"
            target      = "host"
            ignore_case = false
            values = [
              "a.tf-teo-t.xyz",
            ]
          }
          and {
            operator    = "equal"
            target      = "extension"
            ignore_case = false
            values = [
              "png",
            ]
          }
        }
        or {
          and {
            operator    = "notequal"
            target      = "filename"
            ignore_case = false
            values = [
              "test.txt",
            ]
          }
        }
        actions {
          normal_action {
            action = "UpstreamUrlRedirect"
            parameters {
              name = "Type"
              values = [
                "Path",
              ]
            }
            parameters {
              name = "Action"
              values = [
                "addPrefix",
              ]
            }
            parameters {
              name = "Value"
              values = [
                "/www",
              ]
            }
          }
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `rule_name` - (Required, String) The rule name (1 to 255 characters).
* `rules` - (Required, List) Rule items list.
* `status` - (Required, String) Rule status. Values: `enable`: Enabled; `disable`: Disabled.
* `zone_id` - (Required, String, ForceNew) ID of the site.
* `tags` - (Optional, Set: [`String`]) rule tag list.

The `actions` object of `rules` supports the following:

* `code_action` - (Optional, List) Define a code action.
* `normal_action` - (Optional, List) Define a normal action.
* `rewrite_action` - (Optional, List) Define a rewrite action.

The `and` object of `or` supports the following:

* `operator` - (Required, String) Condition operator. Valid values are `equal`, `notequal`.
* `target` - (Required, String) Condition target. Valid values:- `host`: Host of the URL.- `filename`: filename of the URL.- `extension`: file extension of the URL.- `full_url`: full url.- `url`: path of the URL.
* `values` - (Required, Set) Condition Value.
* `ignore_case` - (Optional, Bool) Whether to ignore the case of the parameter value, the default value is false.
* `name` - (Optional, String) The parameter name corresponding to the matching type is valid when the Target value is the following, and the valid value cannot be empty: `query_string` (query string): The parameter name of the query string in the URL request under the current site, such as lang and version in lang=cn&version=1; `request_header` (HTTP request header): HTTP request header field name, such as Accept-Language in Accept-Language:zh-CN,zh;q=0.9.

The `and` object of `or` supports the following:

* `operator` - (Required, String) Condition operator. Valid values are `equal`, `notequal`.
* `target` - (Required, String) Condition target. Valid values:- `host`: Host of the URL.- `filename`: filename of the URL.- `extension`: file extension of the URL.- `full_url`: full url.- `url`: path of the URL.
* `values` - (Required, Set) Condition Value.
* `ignore_case` - (Optional, Bool) Whether to ignore the case of the parameter value, the default value is false.
* `name` - (Optional, String) The parameter name corresponding to the matching type is valid when the Target value is the following, and the valid value cannot be empty:- `query_string` (query string): The parameter name of the query string in the URL request under the current site, such as lang and version in lang=cn&version=1; `request_header` (HTTP request header): HTTP request header field name, such as Accept-Language in Accept-Language:zh-CN,zh;q=0.9.

The `code_action` object of `actions` supports the following:

* `action` - (Required, String) Action name.
* `parameters` - (Required, List) Action parameters.

The `normal_action` object of `actions` supports the following:

* `action` - (Required, String) Action name.
* `parameters` - (Required, List) Action parameters.

The `or` object of `rules` supports the following:

* `and` - (Required, List) AND Conditions list of the rule. Rule would be triggered if all conditions are true.

The `parameters` object of `code_action` supports the following:

* `name` - (Required, String) Parameter Name.
* `status_code` - (Required, Int) HTTP status code to use.
* `values` - (Required, Set) Parameter Values.

The `parameters` object of `normal_action` supports the following:

* `name` - (Required, String) Parameter Name.
* `values` - (Required, Set) Parameter Values.

The `parameters` object of `rewrite_action` supports the following:

* `action` - (Required, String) Action to take on the HEADER. Valid values: `add`, `del`, `set`.
* `name` - (Required, String) Target HEADER name.
* `values` - (Required, Set) Parameter Value.

The `rewrite_action` object of `actions` supports the following:

* `action` - (Required, String) Action name.
* `parameters` - (Required, List) Action parameters.

The `rules` object of `sub_rules` supports the following:

* `actions` - (Required, List) Actions list of the rule. See details in data source `rule_engine_setting`.
* `or` - (Required, List) OR Conditions list of the rule. Rule would be triggered if any of the condition is true.

The `rules` object supports the following:

* `actions` - (Required, List) Actions list of the rule. See details in data source `rule_engine_setting`.
* `or` - (Required, List) OR Conditions list of the rule. Rule would be triggered if any of the condition is true.
* `sub_rules` - (Optional, List) Actions list of the rule. See details in data source `rule_engine_setting`.

The `sub_rules` object of `rules` supports the following:

* `rules` - (Required, List) Rule items list.
* `tags` - (Optional, Set) rule tag list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


## Import

teo rule_engine can be imported using the id#rule_id, e.g.
```
terraform import tencentcloud_teo_rule_engine.rule_engine zone-297z8rf93cfw#rule-ajol584a
```

