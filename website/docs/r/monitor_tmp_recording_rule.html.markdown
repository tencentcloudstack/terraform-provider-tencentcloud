---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_tmp_recording_rule"
sidebar_current: "docs-tencentcloud-resource-monitor_tmp_recording_rule"
description: |-
  Provides a resource to create a monitor tmp recordingRule
---

# tencentcloud_monitor_tmp_recording_rule

Provides a resource to create a monitor tmp recordingRule

## Example Usage

```hcl
resource "tencentcloud_monitor_tmp_recording_rule" "recordingRule" {
  name        = "dasdasdsadasd"
  group       = <<EOF
---
name: example-test
rules:
  - record: job:http_inprogress_requests:sum
    expr: sum by (job) (http_inprogress_requests)
EOF
  instance_id = "prom-c89b3b3u"
  rule_state  = 2
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, String) Recording rule group.
* `instance_id` - (Required, String) Instance id.
* `name` - (Required, String) Recording rule name.
* `rule_state` - (Optional, Int) Rule state.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

monitor recordingRule can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_recording_rule.recordingRule instanceId#recordingRule_id
```

