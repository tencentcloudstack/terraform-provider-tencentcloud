---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_trigger_config"
sidebar_current: "docs-tencentcloud-resource-scf_trigger_config"
description: |-
  Provides a resource to create a scf trigger_config
---

# tencentcloud_scf_trigger_config

Provides a resource to create a scf trigger_config

## Example Usage

```hcl
resource "tencentcloud_scf_trigger_config" "trigger_config" {
  enable          = "OPEN"
  function_name   = "keep-1676351130"
  trigger_name    = "SCF-timer-1685540160"
  type            = "timer"
  qualifier       = "$DEFAULT"
  namespace       = "default"
  trigger_desc    = "* 1 2 * * * *"
  description     = "func"
  custom_argument = "Information"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String, ForceNew) function name.
* `trigger_name` - (Required, String, ForceNew) trigger Name.
* `type` - (Required, String) trigger type.
* `custom_argument` - (Optional, String) User Additional Information.
* `description` - (Optional, String) Trigger description.
* `enable` - (Optional, String) The trigger is turned on or off, and the parameter passed as OPEN is turned on, and CLOSE is turned off.
* `namespace` - (Optional, String, ForceNew) function namespace.
* `qualifier` - (Optional, String) Function version. It defaults to `$LATEST`. It's recommended to use `[$DEFAULT](https://intl.cloud.tencent.com/document/product/583/36149?from_cn_redirect=1#.E9.BB.98.E8.AE.A4.E5.88.AB.E5.90.8D)` for canary release.
* `trigger_desc` - (Optional, String) TriggerDesc parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

scf trigger_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_trigger_config.trigger_config functionName#namespace#triggerName
```

