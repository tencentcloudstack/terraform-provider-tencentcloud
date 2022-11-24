---
subcategory: "TencentCloud Automation Tools(TAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tat_invoker"
sidebar_current: "docs-tencentcloud-resource-tat_invoker"
description: |-
  Provides a resource to create a tat invoker
---

# tencentcloud_tat_invoker

Provides a resource to create a tat invoker

## Example Usage

```hcl
resource "tencentcloud_tat_invoker" "invoker" {
  name         = "pwd-1"
  type         = "SCHEDULE"
  command_id   = "cmd-6fydo27j"
  instance_ids = ["ins-3c7q2ebs", ]
  username     = "root"
  # parameters = ""
  schedule_settings {
    policy = "ONCE"
    # recurrence = ""
    invoke_time = "2099-11-17T16:00:00Z"
  }
}
```

## Argument Reference

The following arguments are supported:

* `command_id` - (Required, String) Remote command ID.
* `instance_ids` - (Required, Set: [`String`]) ID of the instance bound to the trigger. Up to 100 IDs are allowed.
* `name` - (Required, String) Invoker name.
* `type` - (Required, String) Invoker type. It can only be `SCHEDULE` (recurring invokers).
* `parameters` - (Optional, String) Custom parameters of the command.
* `schedule_settings` - (Optional, List) Settings required for a recurring invoker.
* `username` - (Optional, String) The user who executes the command.

The `schedule_settings` object supports the following:

* `policy` - (Required, String) Execution policy: `ONCE`: Execute once; `RECURRENCE`: Execute repeatedly.
* `invoke_time` - (Optional, String) The next execution time of the invoker. This field is required if Policy is ONCE.
* `recurrence` - (Optional, String) Trigger the crontab expression. This field is required if `Policy` is `RECURRENCE`. The crontab expression is parsed in UTC+8.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Creation time.
* `enable` - Whether to enable the invoker.
* `invoker_id` - Invoker ID.
* `updated_time` - Modification time.


## Import

tat invoker can be imported using the id, e.g.
```
$ terraform import tencentcloud_tat_invoker.invoker ivk-gwb4ztk5
```

