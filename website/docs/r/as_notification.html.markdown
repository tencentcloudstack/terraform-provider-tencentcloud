---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_notification"
sidebar_current: "docs-tencentcloud-resource-as_notification"
description: |-
  Provides a resource for an AS (Auto scaling) notification.
---

# tencentcloud_as_notification

Provides a resource for an AS (Auto scaling) notification.

## Example Usage

```hcl
resource "tencentcloud_autoscaling_notification" "aslab" {
  scaling_group_id              = "sg-12af45"
  notification_type             = ["SCALE_OUT_FAILED", "SCALE_IN_SUCCESSFUL", "SCALE_IN_FAILED", "REPLACE_UNHEALTHY_INSTANCE_FAILED"]
  notification_user_group_ids   = ["ASGID"]
}
```

## Argument Reference

The following arguments are supported:

* `notification_type` - (Required) A list of Notification Types that trigger notifications. Acceptable values are SCALE_OUT_FAILED, SCALE_IN_SUCCESSFUL, SCALE_IN_FAILED, REPLACE_UNHEALTHY_INSTANCE_SUCCESSFUL and REPLACE_UNHEALTHY_INSTANCE_FAILED.
* `notification_user_group_ids` - (Required) A group of user IDs to be notified.
* `scaling_group_id` - (Required, ForceNew) ID of a scaling group.


