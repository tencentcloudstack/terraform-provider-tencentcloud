# Add tencentcloud_config_alarm_policy Resource

## What

Add a new Terraform resource `tencentcloud_config_alarm_policy` for managing Tencent Cloud Config alarm policies. This resource supports creating, reading, updating, and deleting alarm policies that monitor configuration compliance events.

## Why

Config alarm policies allow users to receive notifications when compliance events occur (e.g., non-compliant resources). Currently no Terraform resource exists for this feature.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `AddAlarmPolicy` | Returns `AlarmPolicyId` (uint64) as resource unique ID |
| Read | `ListAlarmPolicy` | Paginated by `Offset`; iterate to find matching `AlarmPolicyId` |
| Update | `UpdateAlarmPolicy` | All fields updatable except ID |
| Delete | `DeleteAlarmPolicy` | Pass `AlarmPolicyId` (uint64) |

## Resource ID

`AlarmPolicyId` (uint64), stored as string via `helper.UInt64ToStr`.
