---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_start_instance_refresh"
sidebar_current: "docs-tencentcloud-resource-as_start_instance_refresh"
description: |-
  Provides a resource to create as instance refresh
---

# tencentcloud_as_start_instance_refresh

Provides a resource to create as instance refresh

## Example Usage

```hcl
resource "tencentcloud_as_start_instance_refresh" "example" {
  auto_scaling_group_id = "asg-8n7fdm28"
  refresh_mode          = "ROLLING_UPDATE_RESET"
  refresh_settings {
    check_instance_target_health = false
    rolling_update_settings {
      batch_number = 1
      batch_pause  = "AUTOMATIC"
      max_surge    = 1
      fail_process = "AUTO_PAUSE"
    }
    check_instance_target_health_timeout = 1800
  }

  timeouts {
    create = "10m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_group_id` - (Required, String, ForceNew) Scaling group ID.
* `refresh_settings` - (Required, List, ForceNew) Refresh settings.
* `refresh_mode` - (Optional, String, ForceNew) Refresh mode. Value range: ROLLING_UPDATE_RESET: Reinstall the system for rolling update; ROLLING_UPDATE_REPLACE: Create a new instance for rolling update. This mode does not support the rollback interface yet.

The `refresh_settings` object supports the following:

* `rolling_update_settings` - (Required, List) Rolling update settings parameters. RefreshMode is the rolling update. This parameter must be filled in.Note: This field may return null, indicating that no valid value can be obtained.
* `check_instance_target_health_timeout` - (Optional, Int) The timeout period for backend service health status checks, in seconds. The valid range is [60, 7200], with a default value of 1800 seconds. This takes effect only when the CheckInstanceTargetHealth parameter is enabled. If the instance health check times out, it will be marked as a refresh failure.
* `check_instance_target_health` - (Optional, Bool) Backend service health check status for instances, defaults to FALSE. This setting takes effect only for scaling groups bound with application load balancers. When enabled, if an instance fails the check after being refreshed, its load balancer port weight remains 0 and is marked as a refresh failure. Valid values: <br><li>TRUE: Enable the check.</li> <li>FALSE: Do not enable the check.

The `rolling_update_settings` object of `refresh_settings` supports the following:

* `batch_number` - (Required, Int) Batch quantity. The batch quantity should be a positive integer greater than 0, but cannot exceed the total number of instances pending refresh.
* `batch_pause` - (Optional, String) Pause policy between batches. Default value: Automatic. Valid values: <br><li>FIRST_BATCH_PAUSE: Pause after the first batch update completes.</li> <li>BATCH_INTERVAL_PAUSE: Pause between each batch update.</li> <li>AUTOMATIC: No pauses.
* `fail_process` - (Optional, String) Failure Handling Policy. The default value is `AUTO_PAUSE`. The values are as follows, `AUTO_PAUSE`: Pause after refresh fails; `AUTO_ROLLBACK`: Roll back after refresh fails; `AUTO_CANCEL`: Cancel after refresh fails.
* `max_surge` - (Optional, Int) Maximum Extra Quantity. After setting this parameter, a batch of pay-as-you-go extra instances will be created according to the launch configuration before the rolling update starts, and the extra instances will be destroyed after the rolling update is completed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.

