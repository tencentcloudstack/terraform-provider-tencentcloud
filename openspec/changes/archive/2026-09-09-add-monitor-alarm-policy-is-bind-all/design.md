## Context

The `tencentcloud_monitor_alarm_policy` resource manages Tencent Cloud Monitor alarm policies. The resource currently supports parameters such as `policy_name`, `monitor_type`, `namespace`, `remark`, `enable`, `conditions`, `event_conditions`, `notice_ids`, `trigger_tasks`, `policy_tag`, `group_by`, `filter`, `hierarchical_notices`, and `notice_content_tmpl_bind_infos`.

**Current state:**
- Resource file: `tencentcloud/services/monitor/resource_tc_monitor_alarm_policy.go`
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724`

**API behavior analysis:**

| API | IsBindAll in Request | IsBindAll in Response |
|-----|----------------------|-----------------------|
| `CreateAlarmPolicy` | Yes (`IsBindAll *int64`, values 0 or 1, default 0) | No |
| `DescribeAlarmPolicy` | N/A | Yes (`Policy.IsBindAll *int64`) |
| `ModifyAlarmPolicyInfo` | No | N/A |
| `ModifyAlarmPolicyStatus` | No | N/A |
| `ModifyAlarmPolicyCondition` | No | N/A |
| `ModifyAlarmPolicyTasks` | No | N/A |
| `ModifyAlarmPolicyNotice` | No | N/A |
| `DeleteAlarmPolicy` | No | N/A |

**Key constraint:** `IsBindAll` is only available in the Create request for writing. No Modify API supports updating `IsBindAll`. The Describe API's `AlarmPolicy` struct includes the `IsBindAll` field, so it can be refreshed on Read.

## Goals / Non-Goals

**Goals:**
- Add `is_bind_all` (Optional, ForceNew, TypeInt) parameter to `tencentcloud_monitor_alarm_policy` with valid values `0` and `1`
- Pass `IsBindAll` to the `CreateAlarmPolicy` API when specified by the user
- Read `IsBindAll` from the `DescribeAlarmPolicy` API response (`response.Policy.IsBindAll`) to support state refresh and import
- Maintain full backward compatibility — existing configurations continue to work unchanged

**Non-Goals:**
- Making `is_bind_all` updatable (no Modify API supports it; must be immutable via ForceNew)
- Adding validation that `is_bind_all` conflicts with `filter` (the API handles this server-side; provider-side enforcement is out of scope)

## Decisions

### Decision 1: `is_bind_all` is ForceNew (immutable after creation)

**Rationale:** No Modify API (`ModifyAlarmPolicyInfo`, `ModifyAlarmPolicyStatus`, `ModifyAlarmPolicyCondition`, `ModifyAlarmPolicyTasks`, `ModifyAlarmPolicyNotice`) accepts an `IsBindAll` parameter. Since the value cannot be updated, the parameter is marked with `ForceNew: true`. Changing `is_bind_all` will trigger resource recreation, which is the standard Terraform pattern for immutable parameters. This is simpler and more idiomatic than the `immutableArgs` array pattern, because there is only a single immutable parameter being added and the existing resource does not use the `immutableArgs` pattern.

### Decision 2: Use `TypeInt` with `ValidateAllowedIntValue` for `is_bind_all`

**Rationale:** The cloud API field `IsBindAll` is `*int64` with valid values `0` (no) and `1` (yes). Using `TypeInt` with `tccommon.ValidateAllowedIntValue([]int{0, 1})` is consistent with the existing `enable` field pattern and the cloud API type.

### Decision 3: Read `IsBindAll` from Describe API with nil check

**Rationale:** The `AlarmPolicy` struct's `IsBindAll` field may return null (per SDK documentation). The Read function must check `policy.IsBindAll != nil` before calling `d.Set("is_bind_all", ...)`, consistent with the existing read patterns in the resource and the project's hard constraint of checking nil before setting fields.

### Decision 4: Set `IsBindAll` in Create only when explicitly provided

**Rationale:** The `is_bind_all` parameter is Optional. In the Create function, read the value from schema data and pass it to the request when the user has set it (using `d.GetOk` or by reading the raw value). Since the API default is `0` (not bind all), only setting it when the user explicitly configures it avoids sending unnecessary parameters and maintains backward compatibility.

## Risks / Trade-offs

- **[Risk] Changing `is_bind_all` destroys and recreates the resource**: Using `ForceNew: true` means changing this parameter triggers resource destruction and recreation.
  - **Mitigation:** This is the expected Terraform behavior for immutable parameters. The resource recreation will delete the old alarm policy and create a new one with the updated `IsBindAll` value. Users are warned via the Terraform plan output that the resource will be replaced.

- **[Risk] `IsBindAll` returns null from Describe API**: The SDK documentation notes `IsBindAll` may return null.
  - **Mitigation:** The Read function checks `policy.IsBindAll != nil` before setting the field, consistent with the existing read patterns in the resource.
