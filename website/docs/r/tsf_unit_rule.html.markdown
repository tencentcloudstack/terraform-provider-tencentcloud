---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_unit_rule"
sidebar_current: "docs-tencentcloud-resource-tsf_unit_rule"
description: |-
  Provides a resource to create a tsf unit_rule
---

# tencentcloud_tsf_unit_rule

Provides a resource to create a tsf unit_rule

## Example Usage

```hcl
resource "tencentcloud_tsf_unit_rule" "unit_rule" {
  gateway_instance_id = "gw-ins-rug79a70"
  name                = "terraform-test"
  description         = "terraform-desc"
  unit_rule_item_list {
    relationship        = "AND"
    dest_namespace_id   = "namespace-y8p88eka"
    dest_namespace_name = "garden-test_default"
    name                = "Rule1"
    description         = "rule1-desc"
    unit_rule_tag_list {
      tag_type     = "U"
      tag_field    = "aaa"
      tag_operator = "IN"
      tag_value    = "1"
    }

  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_instance_id` - (Required, String) gateway entity ID.
* `name` - (Required, String) rule name.
* `description` - (Optional, String) rule description.
* `unit_rule_item_list` - (Optional, List) list of rule items.

The `unit_rule_item_list` object supports the following:

* `dest_namespace_id` - (Required, String) destination namespace ID.
* `dest_namespace_name` - (Required, String) destination namespace name.
* `name` - (Required, String) rule item name.
* `relationship` - (Required, String) logical relationship: AND/OR.
* `description` - (Optional, String) rule description.
* `priority` - (Optional, Int) rule order, the smaller the higher the priority: the default is 0.
* `rule_id` - (Optional, String) rule item ID.
* `unit_rule_id` - (Optional, String) Unitization rule ID.
* `unit_rule_tag_list` - (Optional, List) list of rule labels.

The `unit_rule_tag_list` object of `unit_rule_item_list` supports the following:

* `tag_field` - (Required, String) label name.
* `tag_operator` - (Required, String) Operator: IN/NOT_IN/EQUAL/NOT_EQUAL/REGEX.
* `tag_type` - (Required, String) Tag Type: U(User Tag).
* `tag_value` - (Required, String) tag value.
* `rule_id` - (Optional, String) rule ID.
* `unit_rule_item_id` - (Optional, String) Unitization rule item ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - rule ID.
* `status` - usage status: enabled/disabled.


## Import

tsf unit_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_unit_rule.unit_rule unit-rl-zbywqeca
```

