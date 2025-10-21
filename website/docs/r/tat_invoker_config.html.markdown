---
subcategory: "TencentCloud Automation Tools(TAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tat_invoker_config"
sidebar_current: "docs-tencentcloud-resource-tat_invoker_config"
description: |-
  Provides a resource to create a tat invoker_config
---

# tencentcloud_tat_invoker_config

Provides a resource to create a tat invoker_config

## Example Usage

```hcl
resource "tencentcloud_tat_invoker_config" "invoker_config" {
  invoker_id     = "ivk-cas4upyf"
  invoker_status = "on"
}
```

## Argument Reference

The following arguments are supported:

* `invoker_id` - (Required, String) ID of the invoker to be enabled.
* `invoker_status` - (Required, String) Invoker on and off state, Values: `on`, `off`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tat invoker_config can be imported using the id, e.g.

```
terraform import tencentcloud_tat_invoker_config.invoker_config invoker_config_id
```

