---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_enable_unit_rule"
sidebar_current: "docs-tencentcloud-resource-tsf_enable_unit_rule"
description: |-
  Provides a resource to create a tsf enable_unit_rule
---

# tencentcloud_tsf_enable_unit_rule

Provides a resource to create a tsf enable_unit_rule

## Example Usage

```hcl
resource "tencentcloud_tsf_enable_unit_rule" "enable_unit_rule" {
  rule_id = "unit-rl-is9m4nxz"
  switch  = "enabled"
}
```

## Argument Reference

The following arguments are supported:

* `rule_id` - (Required, String) api ID.
* `switch` - (Required, String) switch, on: `enabled`, off: `disabled`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tsf enable_unit_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_enable_unit_rule.enable_unit_rule enable_unit_rule_id
```

