Provides a resource to create a teo rule_engine

~> **NOTE:** The current resource has been deprecated, please use `tencentcloud_teo_l7_acc_rule`.

Example Usage

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
          name   = "Type"
          values = [
            "Path",
          ]
        }
        parameters {
          name   = "Action"
          values = [
            "addPrefix",
          ]
        }
        parameters {
          name   = "Value"
          values = [
            "/sss",
          ]
        }
      }
    }

    or {
      and {
        operator = "equal"
        target   = "host"
        ignore_case = false
        values   = [
          "a.tf-teo-t.xyz",
        ]
      }
      and {
        operator = "equal"
        target   = "extension"
        ignore_case = false
        values   = [
          "jpg",
        ]
      }
    }
    or {
      and {
        operator = "equal"
        target   = "filename"
        ignore_case = false
        values   = [
          "test.txt",
        ]
      }
    }

    sub_rules {
      tags = ["png"]
      rules {
        or {
          and {
            operator = "notequal"
            target   = "host"
            ignore_case = false
            values   = [
              "a.tf-teo-t.xyz",
            ]
          }
          and {
            operator = "equal"
            target   = "extension"
            ignore_case = false
            values   = [
              "png",
            ]
          }
        }
        or {
          and {
            operator = "notequal"
            target   = "filename"
            ignore_case = false
            values   = [
              "test.txt",
            ]
          }
        }
        actions {
          normal_action {
            action = "UpstreamUrlRedirect"
            parameters {
              name   = "Type"
              values = [
                "Path",
              ]
            }
            parameters {
              name   = "Action"
              values = [
                "addPrefix",
              ]
            }
            parameters {
              name   = "Value"
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

* `zone_id` - (Required, ForceNew) ID of the site.
* `rule_name` - (Required) The rule name (1 to 255 characters).
* `status` - (Required) Rule status. Values: `enable`, `disable`.
* `tags` - (Optional) Rule tag list.
* `rules` - (Required) Rule items list. See [rules block](#rules-block) below.

### rules Block

The `rules` block supports:

* `actions` - (Optional) Features to be executed. See [actions block](#actions-block) below.
* `or` - (Required) OR Conditions list of the rule.
* `sub_rules` - (Optional) Nested rules. See [sub_rules block](#sub_rules-block) below.

### actions Block

The `actions` block supports:

* `normal_action` - (Optional) Common operation.
* `rewrite_action` - (Optional) Feature operation with a request/response header.
* `code_action` - (Optional) Feature operation with a status code.

### sub_rules Block

The `sub_rules` block supports:

* `tags` - (Optional) Tag of the rule.
* `rules` - (Required) Nested rule settings.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `rule_id` - Rule ID.
* `rule_priority` - Rule priority, larger value, higher priority, minimum is 1.
* `rule_items` - (Computed) Rule items returned from DescribeRules API. This field provides the complete rule configuration structure from the API response, including all rule details such as rule type, conditions, and actions.
  * `rule_id` - Rule ID.
  * `rule_name` - The rule name (1 to 255 characters).
  * `status` - Rule status. Values: `enable`, `disable`.
  * `rule_priority` - Rule priority, larger value, higher priority, minimum is 1.
  * `tags` - Rule tag list.
  * `rules` - Rule content, containing conditions, actions, and sub_rules.

Import

teo rule_engine can be imported using the id#rule_id, e.g.
```
terraform import tencentcloud_teo_rule_engine.rule_engine zone-297z8rf93cfw#rule-ajol584a
```