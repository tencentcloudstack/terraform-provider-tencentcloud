## 1. Service Layer

- [x] 1.1 Append `DescribeConfigAlarmPolicyById()` — calls `ListAlarmPolicy` iterating pages until matching `AlarmPolicyId` found

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_config_alarm_policy.go` with schema
- [x] 2.2 Implement Create: call `AddAlarmPolicy`, set ID to string(`AlarmPolicyId`)
- [x] 2.3 Implement Read: call service, populate state from `AlarmPolicyRsp`
- [x] 2.4 Implement Update: call `UpdateAlarmPolicy` on any changed field
- [x] 2.5 Implement Delete: call `DeleteAlarmPolicy`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_config_alarm_policy` in `provider.go` ResourcesMap

## 4. Documentation & Tests

- [x] 4.1 Create `resource_tc_config_alarm_policy.md`
- [x] 4.2 Create `resource_tc_config_alarm_policy_test.go`

