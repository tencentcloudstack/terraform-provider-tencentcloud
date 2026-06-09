---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_start_config_rule_evaluation_operation"
sidebar_current: "docs-tencentcloud-resource-config_start_config_rule_evaluation_operation"
description: |-
  Provides a resource to trigger a Config rule evaluation (one-shot operation).
---

# tencentcloud_config_start_config_rule_evaluation_operation

Provides a resource to trigger a Config rule evaluation (one-shot operation).

## Example Usage

### Trigger evaluation by rule ID

```hcl
resource "tencentcloud_config_start_config_rule_evaluation_operation" "example" {
  rule_id = "cr-xhsd76j603v0a8ma0i73"
}
```

### Trigger evaluation by compliance pack ID

```hcl
resource "tencentcloud_config_start_config_rule_evaluation_operation" "example" {
  compliance_pack_id = "cp-3kr5im1ssbg6tdo5jbi9"
}
```

## Argument Reference

The following arguments are supported:

* `compliance_pack_id` - (Optional, String, ForceNew) Compliance pack ID to trigger evaluation for.
* `rule_id` - (Optional, String, ForceNew) Config rule ID to trigger evaluation for.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



