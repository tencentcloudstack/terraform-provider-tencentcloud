---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_policy_set_default"
sidebar_current: "docs-tencentcloud-resource-monitor_alarm_policy_set_default"
description: |-
  Provides a resource to create a monitor policy_set_default
---

# tencentcloud_monitor_alarm_policy_set_default

Provides a resource to create a monitor policy_set_default

## Example Usage

```hcl
resource "tencentcloud_monitor_alarm_policy_set_default" "policy_set_default" {
  module    = "monitor"
  policy_id = "policy-u4iykjkt"
}
```

## Argument Reference

The following arguments are supported:

* `module` - (Required, String, ForceNew) Fixed value, as `monitor`.
* `policy_id` - (Required, String, ForceNew) Policy id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



