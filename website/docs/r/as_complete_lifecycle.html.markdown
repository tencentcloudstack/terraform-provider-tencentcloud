---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_complete_lifecycle"
sidebar_current: "docs-tencentcloud-resource-as_complete_lifecycle"
description: |-
  Provides a resource to create a as complete_lifecycle
---

# tencentcloud_as_complete_lifecycle

Provides a resource to create a as complete_lifecycle

## Example Usage

```hcl
resource "tencentcloud_as_complete_lifecycle" "complete_lifecycle" {
  lifecycle_hook_id       = "ash-xxxxxxxx"
  lifecycle_action_result = "CONTINUE"
  instance_id             = "ins-xxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `lifecycle_action_result` - (Required, String, ForceNew) Result of the lifecycle action. Value range: `CONTINUE`, `ABANDON`.
* `lifecycle_hook_id` - (Required, String, ForceNew) Lifecycle hook ID.
* `instance_id` - (Optional, String, ForceNew) Instance ID. Either InstanceId or LifecycleActionToken must be specified.
* `lifecycle_action_token` - (Optional, String, ForceNew) Either InstanceId or LifecycleActionToken must be specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



