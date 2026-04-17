## 1. Service Layer

- [x] 1.1 Append `DescribeConfigRuleById()` — wraps `DescribeConfigRule` with Retry

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_config_rule.go` with full schema
- [x] 2.2 Implement Create: call `AddConfigRule`, set ID to `RuleId`; call Open/Close if status provided
- [x] 2.3 Implement Read: call service, populate all fields from `ConfigRule`
- [x] 2.4 Implement Update: call `UpdateConfigRule` for content changes; call Open/Close for status changes
- [x] 2.5 Implement Delete: call `DeleteConfigRule`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_config_rule` in `provider.go` ResourcesMap

## 4. Documentation & Tests

- [x] 4.1 Create `resource_tc_config_rule.md`
- [x] 4.2 Create `resource_tc_config_rule_test.go`

