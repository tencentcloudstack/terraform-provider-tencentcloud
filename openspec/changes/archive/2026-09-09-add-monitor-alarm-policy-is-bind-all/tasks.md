## 1. Resource Schema Changes

- [x] 1.1 Add `is_bind_all` schema field (TypeInt, Optional, ForceNew, ValidateAllowedIntValue [0, 1]) to `ResourceTencentCloudMonitorAlarmPolicy()` in `tencentcloud/services/monitor/resource_tc_monitor_alarm_policy.go`

## 2. Create Function Changes

- [x] 2.1 Read `is_bind_all` from schema data in `resourceTencentMonitorAlarmPolicyCreate` and pass `IsBindAll` to the `CreateAlarmPolicy` API request when the user has set it

## 3. Read Function Changes

- [x] 3.1 Read `policy.IsBindAll` from the `DescribeAlarmPolicy` API response in `resourceTencentMonitorAlarmPolicyRead`
- [x] 3.2 Check `policy.IsBindAll != nil` before calling `d.Set("is_bind_all", ...)` to handle null response fields

## 4. Unit Test Changes

- [x] 4.1 Add unit test cases for the `is_bind_all` parameter in `tencentcloud/services/monitor/resource_tc_monitor_alarm_policy_test.go` using mock (gomonkey) approach for cloud API mocking

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/monitor/resource_tc_monitor_alarm_policy.md` with `is_bind_all` usage example and description

## 6. Validation

- [x] 6.1 Verify the code compiles successfully
- [x] 6.2 Verify no lint errors
