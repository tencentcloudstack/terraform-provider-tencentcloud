# Design: tencentcloud_config_alarm_policy Resource

## Architecture

Follows `tencentcloud_igtm_strategy` style:

```
provider.go
    └─ resource_tc_config_alarm_policy.go (CRUD handlers)
           └─ service_tencentcloud_config.go (DescribeConfigAlarmPolicyById)
                  └─ config SDK v20220802
```

## File Layout

| File | Action |
|---|---|
| `resource_tc_config_alarm_policy.go` | New |
| `resource_tc_config_alarm_policy.md` | New |
| `resource_tc_config_alarm_policy_test.go` | New |
| `service_tencentcloud_config.go` | Modified — append `DescribeConfigAlarmPolicyById` |
| `provider.go` | Modified — register resource |

## Schema

### Required

| Field | Type | Description |
|---|---|---|
| `name` | String | Alarm policy name |
| `event_scope` | List of Int | Event scope: 1 (current account), 2 (multi-account) |
| `risk_level` | List of Int | Risk level: 1 (high), 2 (medium), 3 (low) |
| `notice_time` | String | Notification time range (e.g. `09:30:00~23:30:00`) |
| `notification_mechanism` | String | Notification mechanism |
| `status` | Int | Status: 1 (enabled), 2 (disabled) |
| `notice_period` | List of Int | Notification weekdays: 1-7 (Mon-Sun) |

### Optional

| Field | Type | Description |
|---|---|---|
| `description` | String | Policy description |

### Computed

| Field | Type | Description |
|---|---|---|
| `alarm_policy_id` | String | Alarm policy unique ID |

## Read Strategy

`ListAlarmPolicy` only takes `Offset` (no Limit). Iterate pages until `AlarmPolicyId` matches or `AlarmPolicyList` is empty. Use `Offset++` as page cursor.

## Mutable Fields

All fields except the computed `alarm_policy_id` are mutable via `UpdateAlarmPolicy`.
