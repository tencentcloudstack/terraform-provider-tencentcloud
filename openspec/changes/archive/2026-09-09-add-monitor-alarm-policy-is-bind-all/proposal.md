## Why

The Tencent Cloud Monitor `CreateAlarmPolicy` API supports an `IsBindAll` parameter that determines whether an alarm policy binds to all objects. When set to `1`, the policy binds to all objects and no `filter` or `BindPolicyObject` call is needed. The Terraform resource `tencentcloud_monitor_alarm_policy` does not currently expose this parameter, so users cannot create "bind all objects" alarm policies through Terraform and are forced to use the console or API directly.

## What Changes

- Add `is_bind_all` (Optional, ForceNew, TypeInt) parameter to the `tencentcloud_monitor_alarm_policy` resource to support specifying whether the alarm policy binds to all objects. Valid values: `0` (no, default) and `1` (yes). This parameter is immutable after creation because no Modify API supports updating `IsBindAll`.
- Pass `IsBindAll` to the `CreateAlarmPolicy` API request when the parameter is set by the user.
- Read `IsBindAll` from the `DescribeAlarmPolicy` API response (`response.Policy.IsBindAll`) to support state refresh and import.
- Maintain full backward compatibility — the new parameter is Optional and defaults to not being set.

## Capabilities

### New Capabilities
- `monitor-alarm-policy-is-bind-all`: Enable the `is_bind_all` parameter on the `tencentcloud_monitor_alarm_policy` resource to allow users to specify whether an alarm policy binds to all objects at creation time.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/monitor/resource_tc_monitor_alarm_policy.go` — add `is_bind_all` schema field, wire through Create flow, add Read support
  - `tencentcloud/services/monitor/resource_tc_monitor_alarm_policy_test.go` — add unit test cases for the new `is_bind_all` parameter
  - `tencentcloud/services/monitor/resource_tc_monitor_alarm_policy.md` — update documentation
- **SDK dependency:** `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724` — `CreateAlarmPolicyRequest` already includes the `IsBindAll *int64` field, and `AlarmPolicy` struct (Describe response) already includes `IsBindAll *int64`. No SDK update is needed.
- **Backward compatibility:** fully backward compatible — the new parameter is Optional and defaults to not being set.
- **API constraints:** `IsBindAll` is only accepted by `CreateAlarmPolicy` (no Modify API supports it), so this parameter is immutable after creation and marked as `ForceNew`.
